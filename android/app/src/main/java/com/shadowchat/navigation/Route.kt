package com.shadowchat.navigation

sealed class Route(val path: String) {
    data object Splash : Route("splash")
    data object Onboarding : Route("onboarding")
    data object CreateIdentity : Route("create_identity")
    data object RecoveryPhrase : Route("recovery_phrase")
    data object Home : Route("home")
    data object ChatList : Route("chat_list")
    data object ChatScreen : Route("chat_screen/{chatId}") {
        fun create(chatId: String) = "chat_screen/$chatId"
    }
    data object NewChat : Route("new_chat")
    data object Contacts : Route("contacts")
    data object QRScanner : Route("qr_scanner")
    data object Settings : Route("settings")
    data object Privacy : Route("privacy")
    data object Security : Route("security")
    data object Profile : Route("profile")
    data object MediaGallery : Route("media_gallery")
}
