package com.shadowchat.core.crypto

import android.security.keystore.KeyGenParameterSpec
import android.security.keystore.KeyProperties
import java.security.KeyStore
import java.security.SecureRandom
import javax.crypto.Cipher
import javax.crypto.KeyGenerator
import javax.crypto.SecretKey
import javax.crypto.spec.GCMParameterSpec
import android.util.Base64

/**
 * Crypto utilities for end-to-end encryption
 */
object CryptoManager {

    private const val KEYSTORE_PROVIDER = "AndroidKeyStore"
    private const val KEY_ALIAS = "shadowchat_master_key"
    private const val TRANSFORMATION = "AES/GCM/NoPadding"
    private const val GCM_TAG_LENGTH = 128
    private const val GCM_IV_LENGTH = 12

    private val keyStore: KeyStore by lazy {
        KeyStore.getInstance(KEYSTORE_PROVIDER).apply { load(null) }
    }

    /**
     * Generate or retrieve the master key from Android Keystore
     */
    fun getOrCreateMasterKey(): SecretKey {
        val existingKey = keyStore.getEntry(KEY_ALIAS, null) as? KeyStore.SecretKeyEntry
        if (existingKey != null) {
            return existingKey.secretKey
        }

        val keyGenerator = KeyGenerator.getInstance(
            KeyProperties.KEY_ALGORITHM_AES,
            KEYSTORE_PROVIDER
        )

        val keyGenSpec = KeyGenParameterSpec.Builder(
            KEY_ALIAS,
            KeyProperties.PURPOSE_ENCRYPT or KeyProperties.PURPOSE_DECRYPT
        )
            .setBlockModes(KeyProperties.BLOCK_MODE_GCM)
            .setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_NONE)
            .setKeySize(256)
            .setUserAuthenticationRequired(false) // Set to true for biometric
            .build()

        keyGenerator.init(keyGenSpec)
        return keyGenerator.generateKey()
    }

    /**
     * Encrypt data using AES-GCM
     */
    fun encrypt(plaintext: ByteArray, key: SecretKey = getOrCreateMasterKey()): EncryptedData {
        val cipher = Cipher.getInstance(TRANSFORMATION)
        cipher.init(Cipher.ENCRYPT_MODE, key)

        val iv = cipher.iv
        val ciphertext = cipher.doFinal(plaintext)

        return EncryptedData(
            ciphertext = ciphertext,
            iv = iv
        )
    }

    /**
     * Decrypt data using AES-GCM
     */
    fun decrypt(encryptedData: EncryptedData, key: SecretKey = getOrCreateMasterKey()): ByteArray {
        val cipher = Cipher.getInstance(TRANSFORMATION)
        val spec = GCMParameterSpec(GCM_TAG_LENGTH, encryptedData.iv)
        cipher.init(Cipher.DECRYPT_MODE, key, spec)

        return cipher.doFinal(encryptedData.ciphertext)
    }

    /**
     * Generate a random passphrase for SQLCipher
     */
    fun generateDatabasePassphrase(): ByteArray {
        val secureRandom = SecureRandom()
        val passphrase = ByteArray(32)
        secureRandom.nextBytes(passphrase)
        return passphrase
    }

    /**
     * Store the database passphrase securely in Keystore
     */
    fun storePassphrase(passphrase: ByteArray): EncryptedData {
        return encrypt(passphrase)
    }

    /**
     * Retrieve the database passphrase from Keystore
     */
    fun retrievePassphrase(encryptedData: EncryptedData): ByteArray {
        return decrypt(encryptedData)
    }

    data class EncryptedData(
        val ciphertext: ByteArray,
        val iv: ByteArray
    ) {
        fun toBase64(): String {
            return Base64.encodeToString(ciphertext, Base64.NO_WRAP) + ":" +
                    Base64.encodeToString(iv, Base64.NO_WRAP)
        }

        companion object {
            fun fromBase64(base64String: String): EncryptedData {
                val parts = base64String.split(":")
                return EncryptedData(
                    ciphertext = Base64.decode(parts[0], Base64.NO_WRAP),
                    iv = Base64.decode(parts[1], Base64.NO_WRAP)
                )
            }
        }
    }
}
