package com.shadowchat

import android.app.Application
import dagger.hilt.android.HiltAndroidApp

@HiltAndroidApp
class ShadowChatApplication : Application() {
    override fun onCreate() {
        super.onCreate()
    }
}
