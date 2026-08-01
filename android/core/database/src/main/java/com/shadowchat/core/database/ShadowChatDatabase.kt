package com.shadowchat.core.database

import android.content.Context
import androidx.room.Database
import androidx.room.Room
import androidx.room.RoomDatabase
import com.shadowchat.core.database.dao.ChatDao
import com.shadowchat.core.database.dao.ContactDao
import com.shadowchat.core.database.dao.MessageDao
import com.shadowchat.core.database.dao.SessionDao
import com.shadowchat.core.database.dao.UserDao
import com.shadowchat.core.database.entity.ChatEntity
import com.shadowchat.core.database.entity.ContactEntity
import com.shadowchat.core.database.entity.MessageEntity
import com.shadowchat.core.database.entity.SessionEntity
import com.shadowchat.core.database.entity.UserEntity
import net.sqlcipher.database.SupportFactory

@Database(
    entities = [
        UserEntity::class,
        ChatEntity::class,
        MessageEntity::class,
        ContactEntity::class,
        SessionEntity::class
    ],
    version = 1,
    exportSchema = false
)
abstract class ShadowChatDatabase : RoomDatabase() {
    abstract fun userDao(): UserDao
    abstract fun chatDao(): ChatDao
    abstract fun messageDao(): MessageDao
    abstract fun contactDao(): ContactDao
    abstract fun sessionDao(): SessionDao

    companion object {
        private const val DATABASE_NAME = "shadowchat.db"

        @Volatile
        private var INSTANCE: ShadowChatDatabase? = null

        fun getInstance(context: Context, passphrase: ByteArray): ShadowChatDatabase {
            return INSTANCE ?: synchronized(this) {
                val factory = SupportFactory(passphrase)
                val instance = Room.databaseBuilder(
                    context.applicationContext,
                    ShadowChatDatabase::class.java,
                    DATABASE_NAME
                )
                    .openHelperFactory(factory)
                    .fallbackToDestructiveMigration()
                    .build()
                INSTANCE = instance
                instance
            }
        }
    }
}
