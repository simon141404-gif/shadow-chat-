package com.shadowchat.core.network.api

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.PATCH
import retrofit2.http.DELETE
import retrofit2.http.Path
import retrofit2.http.Query

interface ShadowChatApi {
    
    // Auth
    @POST("v1/auth/anonymous")
    suspend fun createAnonymousIdentity(): AnonymousResponse

    @POST("v1/auth/refresh")
    suspend fun refreshSession(@Body request: RefreshRequest): TokenResponse

    // Chats
    @GET("v1/chats")
    suspend fun getChats(): List<ChatResponse>

    @POST("v1/chats")
    suspend fun createChat(@Body request: CreateChatRequest): ChatResponse

    @GET("v1/chats/{chatId}")
    suspend fun getChat(@Path("chatId") chatId: String): ChatResponse

    @GET("v1/chats/{chatId}/messages")
    suspend fun getMessages(
        @Path("chatId") chatId: String,
        @Query("limit") limit: Int = 50,
        @Query("offset") offset: Int = 0
    ): List<MessageResponse>

    @POST("v1/chats/{chatId}/messages")
    suspend fun sendMessage(
        @Path("chatId") chatId: String,
        @Body request: SendMessageRequest
    ): MessageResponse

    // Messages
    @PATCH("v1/messages/{messageId}")
    suspend fun editMessage(
        @Path("messageId") messageId: String,
        @Body request: EditMessageRequest
    ): MessageResponse

    @DELETE("v1/messages/{messageId}")
    suspend fun deleteMessage(@Path("messageId") messageId: String)

    // Contacts
    @GET("v1/contacts")
    suspend fun getContacts(): List<ContactResponse>

    @POST("v1/contacts/share")
    suspend fun shareContact(@Body request: ShareContactRequest): ContactResponse

    // Profile
    @GET("v1/profile")
    suspend fun getProfile(): ProfileResponse

    @PATCH("v1/profile")
    suspend fun updateProfile(@Body request: UpdateProfileRequest): ProfileResponse
}

// Request/Response models
@JsonClass(generateAdapter = true)
data class AnonymousResponse(
    @Json(name = "userId") val userId: String,
    @Json(name = "publicId") val publicId: String,
    @Json(name = "token") val token: String
)

@JsonClass(generateAdapter = true)
data class RefreshRequest(
    @Json(name = "refreshToken") val refreshToken: String
)

@JsonClass(generateAdapter = true)
data class TokenResponse(
    @Json(name = "token") val token: String,
    @Json(name = "jti") val jti: String
)

@JsonClass(generateAdapter = true)
data class ChatResponse(
    @Json(name = "id") val id: String,
    @Json(name = "type") val type: String,
    @Json(name = "name") val name: String?,
    @Json(name = "avatarUrl") val avatarUrl: String?,
    @Json(name = "createdAt") val createdAt: String,
    @Json(name = "updatedAt") val updatedAt: String
)

@JsonClass(generateAdapter = true)
data class CreateChatRequest(
    @Json(name = "type") val type: String,
    @Json(name = "name") val name: String?,
    @Json(name = "members") val members: List<String>?
)

@JsonClass(generateAdapter = true)
data class MessageResponse(
    @Json(name = "id") val id: String,
    @Json(name = "chatId") val chatId: String,
    @Json(name = "senderUserId") val senderUserId: String,
    @Json(name = "clientMsgId") val clientMsgId: String,
    @Json(name = "messageType") val messageType: String,
    @Json(name = "content") val content: String,
    @Json(name = "createdAt") val createdAt: String,
    @Json(name = "editedAt") val editedAt: String?
)

@JsonClass(generateAdapter = true)
data class SendMessageRequest(
    @Json(name = "clientMsgId") val clientMsgId: String,
    @Json(name = "messageType") val messageType: String,
    @Json(name = "content") val content: String,
    @Json(name = "replyToMessageId") val replyToMessageId: String?
)

@JsonClass(generateAdapter = true)
data class EditMessageRequest(
    @Json(name = "content") val content: String
)

@JsonClass(generateAdapter = true)
data class ContactResponse(
    @Json(name = "id") val id: String,
    @Json(name = "ownerUserId") val ownerUserId: String,
    @Json(name = "contactUserId") val contactUserId: String,
    @Json(name = "displayName") val displayName: String?
)

@JsonClass(generateAdapter = true)
data class ShareContactRequest(
    @Json(name = "publicId") val publicId: String
)

@JsonClass(generateAdapter = true)
data class ProfileResponse(
    @Json(name = "id") val id: String,
    @Json(name = "publicId") val publicId: String,
    @Json(name = "displayName") val displayName: String?,
    @Json(name = "avatarUrl") val avatarUrl: String?,
    @Json(name = "bio") val bio: String?
)

@JsonClass(generateAdapter = true)
data class UpdateProfileRequest(
    @Json(name = "displayName") val displayName: String?,
    @Json(name = "avatarUrl") val avatarUrl: String?,
    @Json(name = "bio") val bio: String?
)
