package com.laioffer.spotify.ui.home

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import coil.compose.AsyncImage
import com.laioffer.spotify.data.model.Album

@Composable
fun HomeScreen(
    state: HomeUiState,
    onAlbumClick: (Int) -> Unit,
    onRetry: () -> Unit,
    modifier: Modifier = Modifier,
) {
    Box(modifier.fillMaxSize()) {
        when {
            state.isLoading && state.sections.isEmpty() -> CircularProgressIndicator(Modifier.align(Alignment.Center))
            state.error != null && state.sections.isEmpty() -> ErrorPanel(state.error, onRetry, Modifier.align(Alignment.Center))
            else -> LazyColumn(
                modifier = Modifier.fillMaxSize(),
                contentPadding = PaddingValues(vertical = 20.dp),
                verticalArrangement = Arrangement.spacedBy(26.dp),
            ) {
                item {
                    Column(Modifier.padding(horizontal = 20.dp)) {
                        Text("Good listening", style = MaterialTheme.typography.headlineMedium, fontWeight = FontWeight.Bold)
                        Text("No account required — every track comes from your local server.", color = MaterialTheme.colorScheme.onSurfaceVariant)
                    }
                }
                items(state.sections, key = { it.sectionTitle }) { section ->
                    Column(verticalArrangement = Arrangement.spacedBy(12.dp)) {
                        Text(
                            section.sectionTitle,
                            modifier = Modifier.padding(horizontal = 20.dp),
                            style = MaterialTheme.typography.titleLarge,
                            fontWeight = FontWeight.Bold,
                        )
                        LazyRow(
                            contentPadding = PaddingValues(horizontal = 20.dp),
                            horizontalArrangement = Arrangement.spacedBy(14.dp),
                        ) {
                            items(section.albums, key = Album::id) { album ->
                                AlbumCard(album, onAlbumClick)
                            }
                        }
                    }
                }
                if (state.error != null) item { ErrorPanel(state.error, onRetry, Modifier.fillMaxWidth()) }
            }
        }
    }
}

@Composable
private fun AlbumCard(album: Album, onClick: (Int) -> Unit) {
    Column(
        Modifier
            .size(width = 164.dp, height = 224.dp)
            .clickable { onClick(album.id) },
        verticalArrangement = Arrangement.spacedBy(6.dp),
    ) {
        AsyncImage(
            model = album.cover,
            contentDescription = "${album.name} cover",
            modifier = Modifier
                .fillMaxWidth()
                .height(164.dp),
            contentScale = ContentScale.Crop,
        )
        Text(album.name, maxLines = 1, overflow = TextOverflow.Ellipsis, fontWeight = FontWeight.SemiBold)
        Text(album.artists, maxLines = 1, overflow = TextOverflow.Ellipsis, color = MaterialTheme.colorScheme.onSurfaceVariant)
    }
}

@Composable
private fun ErrorPanel(message: String, onRetry: () -> Unit, modifier: Modifier = Modifier) {
    Column(modifier.padding(24.dp), horizontalAlignment = Alignment.CenterHorizontally) {
        Text(message, color = MaterialTheme.colorScheme.error)
        Button(onClick = onRetry, modifier = Modifier.padding(top = 12.dp)) { Text("Try again") }
    }
}
