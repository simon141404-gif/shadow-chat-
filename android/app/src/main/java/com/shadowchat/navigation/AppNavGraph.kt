package com.shadowchat.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.NavHostController
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.navArgument
import com.shadowchat.presentation.chat.ChatScreen
import com.shadowchat.presentation.home.HomeScreen
import com.shadowchat.presentation.onboarding.OnboardingScreen
import com.shadowchat.presentation.qr.QRScannerScreen
import com.shadowchat.presentation.settings.SettingsScreen
import com.shadowchat.presentation.splash.SplashScreen

@Composable
fun AppNavGraph(
    navController: NavHostController = rememberNavController()
) {
    NavHost(
        navController = navController,
        startDestination = Route.Splash.path
    ) {
        composable(Route.Splash.path) {
            SplashScreen(
                onDone = {
                    navController.navigate(Route.Onboarding.path) {
                        popUpTo(Route.Splash.path) { inclusive = true }
                    }
                }
            )
        }

        composable(Route.Onboarding.path) {
            OnboardingScreen(
                onContinue = {
                    navController.navigate(Route.Home.path) {
                        popUpTo(Route.Onboarding.path) { inclusive = true }
                    }
                }
            )
        }

        composable(Route.Home.path) {
            HomeScreen(
                onOpenChat = { chatId ->
                    navController.navigate(Route.ChatScreen.create(chatId))
                },
                onQr = {
                    navController.navigate(Route.QRScanner.path)
                },
                onSettings = {
                    navController.navigate(Route.Settings.path)
                }
            )
        }

        composable(Route.QRScanner.path) {
            QRScannerScreen(
                onBack = { navController.popBackStack() },
                onCodeScanned = { code ->
                    // Handle scanned QR code
                    navController.popBackStack()
                }
            )
        }

        composable(
            route = Route.ChatScreen.path,
            arguments = listOf(navArgument("chatId") { type = NavType.StringType })
        ) { backStackEntry ->
            val chatId = backStackEntry.arguments?.getString("chatId") ?: ""
            ChatScreen(
                chatId = chatId,
                onBack = { navController.popBackStack() }
            )
        }

        composable(Route.Settings.path) {
            SettingsScreen(
                onBack = { navController.popBackStack() }
            )
        }
    }
}
