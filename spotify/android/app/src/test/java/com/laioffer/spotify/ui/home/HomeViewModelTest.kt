package com.laioffer.spotify.ui.home

import com.laioffer.spotify.MainDispatcherRule
import com.laioffer.spotify.data.model.Album
import com.laioffer.spotify.data.model.Section
import com.laioffer.spotify.data.repository.FeedRepository
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.runTest
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertNull
import org.junit.Rule
import org.junit.Test

@OptIn(ExperimentalCoroutinesApi::class)
class HomeViewModelTest {
    @get:Rule val mainDispatcherRule = MainDispatcherRule()

    @Test
    fun `initial load publishes feed content`() = runTest {
        val section = Section("Local", listOf(album))
        val viewModel = HomeViewModel(FakeFeedRepository(result = Result.success(listOf(section))))

        advanceUntilIdle()

        assertFalse(viewModel.uiState.value.isLoading)
        assertNull(viewModel.uiState.value.error)
        assertEquals(section, viewModel.uiState.value.sections.single())
    }

    @Test
    fun `failed load exposes retryable error`() = runTest {
        val viewModel = HomeViewModel(FakeFeedRepository(Result.failure(IllegalStateException("offline"))))

        advanceUntilIdle()

        assertFalse(viewModel.uiState.value.isLoading)
        assertEquals("offline", viewModel.uiState.value.error)
    }

    private class FakeFeedRepository(private val result: Result<List<Section>>) : FeedRepository {
        override suspend fun feed(forceRefresh: Boolean): List<Section> = result.getOrThrow()
        override suspend fun album(id: Int): Album = result.getOrThrow().flatMap { it.albums }.first { it.id == id }
    }

    private companion object {
        val album = Album(1, "Album", "2026", "cover", "Artist", "Description")
    }
}
