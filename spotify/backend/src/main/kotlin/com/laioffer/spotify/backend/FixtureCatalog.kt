package com.laioffer.spotify.backend

import kotlinx.serialization.builtins.ListSerializer
import kotlinx.serialization.json.Json
import java.net.URI

class FixtureCatalog private constructor(
    val feed: List<Section>,
    val playlists: List<Playlist>,
) {
    private val playlistsById = playlists.associateBy(Playlist::id)
    private val audioFixtures = playlists.flatMap(Playlist::songs)
        .map { URI.create(it.src).path.substringAfterLast('/') }
        .toSet()

    fun playlist(id: Int): Playlist? = playlistsById[id]
    fun hasAudioFixture(file: String): Boolean = file in audioFixtures

    companion object {
        private val json = Json { ignoreUnknownKeys = false }

        fun load(
            classLoader: ClassLoader = FixtureCatalog::class.java.classLoader,
            publicBaseUrl: String = System.getenv("PUBLIC_BASE_URL") ?: "http://10.0.2.2:8080",
        ): FixtureCatalog {
            val baseUrl = validateBaseUrl(publicBaseUrl)
            val rawFeed = json.decodeFromString(
                ListSerializer(Section.serializer()),
                classLoader.resourceText("feed.json"),
            )
            val rawPlaylists = json.decodeFromString(
                ListSerializer(Playlist.serializer()),
                classLoader.resourceText("playlists.json"),
            )

            return FixtureCatalog(
                feed = rawFeed.map { section ->
                    section.copy(albums = section.albums.map { album ->
                        album.copy(cover = album.cover.absoluteAgainst(baseUrl))
                    })
                },
                playlists = rawPlaylists.map { playlist ->
                    playlist.copy(songs = playlist.songs.map { song ->
                        song.copy(src = song.src.absoluteAgainst(baseUrl))
                    })
                },
            ).also { it.validate(baseUrl) }
        }

        private fun FixtureCatalog.validate(baseUrl: String) {
            val albumIds = feed.flatMap(Section::albums).map(Album::id)
            require(albumIds.isNotEmpty()) { "feed.json must contain at least one album" }
            require(albumIds.all { it > 0 }) { "Album ids must be positive" }
            require(albumIds.size == albumIds.distinct().size) { "Album ids must be unique" }
            require(playlists.map(Playlist::id).distinct().size == playlists.size) {
                "Playlist ids must be unique"
            }
            require(playlists.all { it.id > 0 }) { "Playlist ids must be positive" }
            require(albumIds.toSet() == playlists.map(Playlist::id).toSet()) {
                "Every album must have exactly one playlist and vice versa"
            }
            require(playlists.all { it.songs.isNotEmpty() }) { "Playlists cannot be empty" }
            val base = URI.create(baseUrl)
            feed.flatMap(Section::albums).forEach { album ->
                val cover = URI.create(album.cover)
                require(cover.sameOriginAs(base) && cover.path == "/covers/${album.id}.svg") {
                    "Album ${album.id} must use its local cover fixture"
                }
            }
            val songUris = playlists.flatMap(Playlist::songs).map { song ->
                require(song.length == "00:05") { "Generated song duration must be 00:05" }
                URI.create(song.src).also { uri ->
                    require(uri.sameOriginAs(base) && Regex("/songs/[a-z0-9-]+\\.wav").matches(uri.path)) {
                        "Songs must use local WAV fixture routes"
                    }
                }
            }
            require(songUris.map(URI::getPath).distinct().size == songUris.size) {
                "Song fixture paths must be unique"
            }
        }

        private fun validateBaseUrl(value: String): String {
            val normalized = value.trim().trimEnd('/')
            val uri = runCatching { URI(normalized) }.getOrNull()
            require(uri != null && uri.scheme?.lowercase() in setOf("http", "https") && !uri.host.isNullOrBlank()) {
                "PUBLIC_BASE_URL must be an absolute HTTP(S) origin"
            }
            require(uri.userInfo == null && uri.query == null && uri.fragment == null && uri.path.isNullOrEmpty()) {
                "PUBLIC_BASE_URL must not include credentials, a path, query, or fragment"
            }
            return normalized
        }

        private fun URI.sameOriginAs(other: URI): Boolean =
            scheme.equals(other.scheme, ignoreCase = true) &&
                host.equals(other.host, ignoreCase = true) &&
                effectivePort() == other.effectivePort()

        private fun URI.effectivePort(): Int = if (port >= 0) port else if (scheme.equals("https", true)) 443 else 80

        private fun ClassLoader.resourceText(name: String): String =
            requireNotNull(getResource(name)) { "Missing classpath resource: $name" }.readText()

        private fun String.absoluteAgainst(baseUrl: String): String = when {
            startsWith("http://") || startsWith("https://") -> this
            startsWith('/') -> "$baseUrl$this"
            else -> "$baseUrl/$this"
        }
    }
}
