package com.laioffer.spotify.backend

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class Album(
    val id: Int,
    @SerialName("album") val name: String,
    val year: String,
    val cover: String,
    val artists: String,
    val description: String,
)

@Serializable
data class Section(
    @SerialName("section_title") val sectionTitle: String,
    val albums: List<Album>,
)

@Serializable
data class Song(
    val name: String,
    val lyric: String,
    val src: String,
    val length: String,
)

@Serializable
data class Playlist(
    val id: Int,
    val songs: List<Song>,
)

@Serializable
data class ApiError(val error: String)

@Serializable
data class HealthResponse(
    val status: String,
    val albums: Int,
    val playlists: Int,
)
