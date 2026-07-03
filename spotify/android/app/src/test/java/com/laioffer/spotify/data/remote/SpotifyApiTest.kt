package com.laioffer.spotify.data.remote

import kotlinx.coroutines.test.runTest
import okhttp3.mockwebserver.MockResponse
import okhttp3.mockwebserver.MockWebServer
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Before
import org.junit.Test
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

class SpotifyApiTest {
    private lateinit var server: MockWebServer
    private lateinit var api: SpotifyApi

    @Before
    fun setUp() {
        server = MockWebServer().apply { start() }
        api = Retrofit.Builder()
            .baseUrl(server.url("/"))
            .addConverterFactory(GsonConverterFactory.create())
            .build()
            .create(SpotifyApi::class.java)
    }

    @After
    fun tearDown() = server.shutdown()

    @Test
    fun `feed maps documented snake case and album field names`() = runTest {
        server.enqueue(
            MockResponse()
                .setHeader("Content-Type", "application/json")
                .setBody(
                    """[{"section_title":"Top mixes","albums":[{"id":7,"album":"Signals","year":"2026","cover":"http://cover","artists":"Fixture","description":"Local"}]}]""",
                ),
        )

        val result = api.getHomeFeed()

        assertEquals("Top mixes", result.single().sectionTitle)
        assertEquals("Signals", result.single().albums.single().name)
        assertEquals("/feed", server.takeRequest().path)
    }

    @Test
    fun `playlist id is inserted into endpoint path`() = runTest {
        server.enqueue(
            MockResponse()
                .setHeader("Content-Type", "application/json")
                .setBody("""{"id":3,"songs":[{"name":"Tone","lyric":"Artist","src":"http://song","length":"00:05"}]}"""),
        )

        val result = api.getPlaylist(3)

        assertEquals(3, result.id)
        assertEquals("Tone", result.songs.single().name)
        assertEquals("/playlist/3", server.takeRequest().path)
    }
}
