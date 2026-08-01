package com.shadowchat.core.database.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "users")
data class UserEntity(
    @PrimaryKey val id: String,
    val publicId: String,
    val displayName: String?,
    val avatarUrl: String?,
    val bio: String?,
    val createdAt: Long,
    val updatedAt: Long
)

@Entity(tableName = "chats")
data class ChatEntity(
    @PrimaryKey val id: String,
    val type: String,
    val name: String?,
    val avatarUrl: String?,
    val lastMessage: String?,
    val lastMessageTime: Long?,
    val createdAt: Long,
    val updatedAt: Long
)

@Entity(tableName = "messages")
data class MessageEntity(
    @PrimaryKey val id: String,
    val chatId: String,
    val senderUserId: String,
    val clientMsgId: String?,
    val messageType: String,
    val content: String,
    val replyToMessageId: String?,
    val createdAt: Long,
    val editedAt: Long?,
    val isFromMe: Boolean
)

@Entity(tableName = "contacts")
data class ContactEntity(
    @PrimaryKey val id: String,
    val contactUserId: String,
    val displayName: String?,
    val avatarUrl: String?,
    val createdAt: Long
)

@Entity(tableName = "sessions")
data class SessionEntity(
    @PrimaryKey val id: String,
    val userId: String,
    val token: String,
    val refreshToken: String,
    val expiresAt: Long,
    val createdAt: Long
)
