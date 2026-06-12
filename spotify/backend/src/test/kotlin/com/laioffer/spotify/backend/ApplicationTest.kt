package com.laioffer.spotify.backend

import io.ktor.client.call.body
import io.ktor.client.request.get
import io.ktor.client.request.header
import io.ktor.client.statement.bodyAsText
import io.ktor.http.ContentType
import io.ktor.http.HttpStatusCode
import io.ktor.http.contentType
import io.ktor.http.HttpHeaders
import io.ktor.server.testing.testApplication
import kotlinx.serialization.builtins.ListSerializer
import kotlinx.serialization.json.Json
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue

class ApplicationTest {
    private val json = Json { ignoreUnknownKeys = true }

    @Test
    fun `feed and playlist endpoints preserve the documented shape`() = testApplication {
        application { spotifyModule(FixtureCatalog.load(publicBaseUrl = "http://example.test")) }

        val feedResponse = client.get("/feed")
        assertEquals(HttpStatusCode.OK, feedResponse.status)
        assertEquals("nosniff", feedResponse.headers["X-Content-Type-Options"])
        assertEquals(ContentType.Application.Json, feedResponse.contentType()?.withoutParameters())
        val feed = json.decodeFromString(ListSerializer(Section.serializer()), feedResponse.bodyAsText())
        assertEquals("Made for local development", feed.first().sectionTitle)
        assertEquals("Midnight Drive", feed.first().albums.first().name)
        assertEquals("http://example.test/covers/1.svg", feed.first().albums.first().cover)

        val playlistResponse = client.get("/playlist/1")
        assertEquals(HttpStatusCode.OK, playlistResponse.status)
        val playlist = json.decodeFromString(Playlist.serializer(), playlistResponse.bodyAsText())
        assertEquals(1, playlist.id)
        assertEquals(3, playlist.songs.size)
        assertTrue(playlist.songs.first().src.startsWith("http://example.test/songs/"))
    }

    @Test
    fun `invalid and unknown playlist ids return useful errors`() = testApplication {
        application { spotifyModule(FixtureCatalog.load()) }

        assertEquals(HttpStatusCode.BadRequest, client.get("/playlist/nope").status)
        assertEquals(HttpStatusCode.NotFound, client.get("/playlist/999").status)
    }

    @Test
    fun `local audio and cover fixtures are playable resources`() = testApplication {
        application { spotifyModule(FixtureCatalog.load()) }

        val audio = client.get("/songs/night-signals.wav")
        assertEquals(HttpStatusCode.OK, audio.status)
        assertEquals("audio/wav", audio.contentType()?.withoutParameters().toString())
        assertEquals("RIFF", audio.body<ByteArray>().take(4).toByteArray().toString(Charsets.US_ASCII))

        val cover = client.get("/covers/1.svg")
        assertEquals(HttpStatusCode.OK, cover.status)
        assertTrue(cover.bodyAsText().contains("<svg"))

        assertEquals(HttpStatusCode.NotFound, client.get("/songs/not-in-the-catalog.wav").status)
        assertEquals(HttpStatusCode.NotFound, client.get("/covers/1").status)

        val range = client.get("/songs/night-signals.wav") {
            header(HttpHeaders.Range, "bytes=0-3")
        }
        assertEquals(HttpStatusCode.PartialContent, range.status)
        assertEquals("RIFF", range.bodyAsText())

        val invalidRange = client.get("/songs/night-signals.wav") {
            header(HttpHeaders.Range, "bytes=999999-1000000")
        }
        assertEquals(HttpStatusCode.RequestedRangeNotSatisfiable, invalidRange.status)
    }

    @Test
    fun `health reports fixture counts`() = testApplication {
        application { spotifyModule(FixtureCatalog.load()) }
        val response = client.get("/health")
        assertEquals(HttpStatusCode.OK, response.status)
        assertTrue(response.bodyAsText().contains("\"albums\": 5"))
    }

    @Test
    fun `port configuration fails fast when invalid`() {
        assertEquals(8080, parsePort(null))
        assertEquals(9090, parsePort("9090"))
        kotlin.test.assertFailsWith<IllegalArgumentException> { parsePort("invalid") }
        kotlin.test.assertFailsWith<IllegalArgumentException> { parsePort("70000") }
    }
}
