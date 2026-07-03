package com.laioffer.spotify.player

import com.laioffer.spotify.MainDispatcherRule
import com.laioffer.spotify.data.model.Album
import com.laioffer.spotify.data.model.Song
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.runTest
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Rule
import org.junit.Test

@OptIn(ExperimentalCoroutinesApi::class)
class PlayerViewModelTest {
    @get:Rule val mainDispatcherRule = MainDispatcherRule()

    @Test
    fun `load selects media and starts playback`() = runTest {
        val controller = FakePlaybackController()
        val viewModel = PlayerViewModel(controller)

        viewModel.load(song, album)
        advanceUntilIdle()

        assertEquals(song.src, controller.loadedUrl)
        assertTrue(controller.playCalled)
        assertEquals(song, viewModel.uiState.value.song)
        assertEquals(album, viewModel.uiState.value.album)
    }

    @Test
    fun `toggle pauses a playing controller and seek delegates`() = runTest {
        val controller = FakePlaybackController()
        val viewModel = PlayerViewModel(controller)
        controller.mutableState.value = PlaybackSnapshot(isPlaying = true, durationMs = 5_000)
        advanceUntilIdle()

        viewModel.togglePlay()
        viewModel.seekTo(1_250)

        assertTrue(controller.pauseCalled)
        assertEquals(1_250L, controller.seekPosition)
    }

    private class FakePlaybackController : PlaybackController {
        val mutableState = MutableStateFlow(PlaybackSnapshot())
        override val state: StateFlow<PlaybackSnapshot> = mutableState
        var loadedUrl: String? = null
        var playCalled = false
        var pauseCalled = false
        var seekPosition: Long? = null

        override fun load(url: String) { loadedUrl = url }
        override fun play() { playCalled = true }
        override fun pause() { pauseCalled = true }
        override fun seekTo(positionMs: Long) { seekPosition = positionMs }
    }

    private companion object {
        val album = Album(1, "Album", "2026", "cover", "Artist", "Description")
        val song = Song("Song", "Artist", "http://localhost/song.wav", "00:05")
    }
}
