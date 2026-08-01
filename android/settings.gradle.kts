pluginManagement {
    repositories {
        google()
        mavenCentral()
        gradlePluginPortal()
    }
}
dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)
    repositories {
        google()
        mavenCentral()
    }
}

rootProject.name = "ShadowChat"
include(":app")
include(":core:common")
include(":core:ui")
include(":core:network")
include(":core:crypto")
include(":core:database")
include(":core:testing")
