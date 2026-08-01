package com.shadowchat.di

import android.content.Context
import com.shadowchat.core.crypto.CryptoManager
import com.shadowchat.core.database.ShadowChatDatabase
import com.shadowchat.core.database.dao.ChatDao
import com.shadowchat.core.database.dao.ContactDao
import com.shadowchat.core.database.dao.MessageDao
import com.shadowchat.core.database.dao.SessionDao
import com.shadowchat.core.database.dao.UserDao
import com.shadowchat.core.network.api.ShadowChatApi
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object AppModule {

    @Provides
    @Singleton
    fun provideDatabase(@ApplicationContext context: Context): ShadowChatDatabase {
        // In production, retrieve the passphrase from secure storage
        val passphrase = CryptoManager.generateDatabasePassphrase()
        return ShadowChatDatabase.getInstance(context, passphrase)
    }

    @Provides
    @Singleton
    fun provideUserDao(database: ShadowChatDatabase): UserDao = database.userDao()

    @Provides
    @Singleton
    fun provideChatDao(database: ShadowChatDatabase): ChatDao = database.chatDao()

    @Provides
    @Singleton
    fun provideMessageDao(database: ShadowChatDatabase): MessageDao = database.messageDao()

    @Provides
    @Singleton
    fun provideContactDao(database: ShadowChatDatabase): ContactDao = database.contactDao()

    @Provides
    @Singleton
    fun provideSessionDao(database: ShadowChatDatabase): SessionDao = database.sessionDao()
}
