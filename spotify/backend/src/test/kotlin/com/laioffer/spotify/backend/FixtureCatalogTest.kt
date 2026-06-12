package com.laioffer.spotify.backend

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertFailsWith

class FixtureCatalogTest {
    @Test
    fun `every feed album has an indexed playlist`() {
        val catalog = FixtureCatalog.load(publicBaseUrl = "http://localhost:8080/")
        val albums = catalog.feed.flatMap(Section::albums)

        assertEquals(5, albums.size)
        albums.forEach { album ->
            assertNotNull(catalog.playlist(album.id))
            assertEquals("http://localhost:8080/covers/${album.id}.svg", album.cover)
        }
    }

    @Test
    fun `public base URL must be a safe origin`() {
        assertFailsWith<IllegalArgumentException> {
            FixtureCatalog.load(publicBaseUrl = "javascript:alert(1)")
        }
        assertFailsWith<IllegalArgumentException> {
            FixtureCatalog.load(publicBaseUrl = "https://example.test/path?query=true")
        }
    }
}
