package com.laioffer.spotify.ui.theme

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

val SpotifyGreen = Color(0xFF1ED760)
val SpotifyBlack = Color(0xFF090909)
val SpotifySurface = Color(0xFF181818)
val SpotifyRaised = Color(0xFF282828)
val MutedText = Color(0xFFB3B3B3)

private val AppColors = darkColorScheme(
    primary = SpotifyGreen,
    onPrimary = Color.Black,
    background = SpotifyBlack,
    onBackground = Color.White,
    surface = SpotifySurface,
    onSurface = Color.White,
    surfaceVariant = SpotifyRaised,
    onSurfaceVariant = MutedText,
    error = Color(0xFFFF6B6B),
)

@Composable
fun SpotifyTheme(
    @Suppress("UNUSED_PARAMETER") darkTheme: Boolean = isSystemInDarkTheme(),
    content: @Composable () -> Unit,
) {
    MaterialTheme(colorScheme = AppColors, content = content)
}
