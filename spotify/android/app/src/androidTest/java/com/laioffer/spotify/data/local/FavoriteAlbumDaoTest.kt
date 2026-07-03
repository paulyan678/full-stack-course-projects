package com.laioffer.spotify.data.local

import android.content.Context
import androidx.room.Room
import androidx.test.core.app.ApplicationProvider
import androidx.test.ext.junit.runners.AndroidJUnit4
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.test.runTest
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class FavoriteAlbumDaoTest {
    private lateinit var database: SpotifyDatabase
    private lateinit var dao: FavoriteAlbumDao

    @Before
    fun createDatabase() {
        val context = ApplicationProvider.getApplicationContext<Context>()
        database = Room.inMemoryDatabaseBuilder(context, SpotifyDatabase::class.java)
            .allowMainThreadQueries()
            .build()
        dao = database.favoriteAlbumDao()
    }

    @After
    fun closeDatabase() = database.close()

    @Test
    fun favoriteCanBeSavedObservedAndRemoved() = runTest {
        val album = FavoriteAlbumEntity(1, "Midnight", "2026", "cover", "Artist", "Description")

        assertFalse(dao.observeIsFavorite(1).first())
        dao.upsert(album)
        assertTrue(dao.observeIsFavorite(1).first())
        assertEquals(album, dao.observeAll().first().single())
        dao.delete(album)
        assertFalse(dao.observeIsFavorite(1).first())
    }
}
