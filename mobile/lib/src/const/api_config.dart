// API Configuration for ChaseApp
// This file configures the app to use either Firebase or the new Go API

import '../../flavors.dart';

enum ApiBackend {
  firebase, // Original Firebase backend
  goApi,    // New Go monolithic API
}

class ApiConfig {
  // Toggle this to switch between backends
  static const ApiBackend _currentBackend = ApiBackend.firebase;
  static ApiBackend get currentBackend => _currentBackend;

  // API Base URLs
  static const String _firebaseBaseUrl = 'https://us-central1-chaseapp-8459b.cloudfunctions.net';
  static const String _goApiBaseUrl = 'https://api.chaseapp.tv';

  static String get baseUrl {
    switch (currentBackend) {
      case ApiBackend.firebase:
        return _firebaseBaseUrl;
      case ApiBackend.goApi:
        return _goApiBaseUrl;
    }
  }

  // API Endpoints mapping
  static String get chatTokenEndpoint {
    switch (currentBackend) {
      case ApiBackend.firebase:
        return '/GetStreamToken';
      case ApiBackend.goApi:
        return '/api/v1/auth/chat-token';
    }
  }

  static String get unsubscribePushEndpoint {
    switch (currentBackend) {
      case ApiBackend.firebase:
        return '/UnsubscribeTopic';
      case ApiBackend.goApi:
        return '/api/v1/push/unsubscribe';
    }
  }

  static String get subscribePushEndpoint {
    switch (currentBackend) {
      case ApiBackend.firebase:
        return '/SubscribeTopic';
      case ApiBackend.goApi:
        return '/api/v1/push/subscribe';
    }
  }

  static String get chasesEndpoint {
    switch (currentBackend) {
      case ApiBackend.firebase:
        return '/chases'; // Firebase function name
      case ApiBackend.goApi:
        return '/api/v1/chases';
    }
  }

  static String get aircraftEndpoint {
    switch (currentBackend) {
      case ApiBackend.firebase:
        return '/getAircraft'; // Firebase function name
      case ApiBackend.goApi:
        return '/api/v1/aircraft';
    }
  }

  static String get weatherEndpoint {
    switch (currentBackend) {
      case ApiBackend.firebase:
        return '/weatherAPI'; // Firebase function name
      case ApiBackend.goApi:
        return '/api/v1/weather/alerts';
    }
  }

  // Helper method to get full URL for an endpoint
  static String getEndpointUrl(String endpoint) {
    return '$baseUrl$endpoint';
  }

  // Check if we're using the new Go API
  static bool get isUsingGoApi => currentBackend == ApiBackend.goApi;

  // Check if we're using Firebase
  static bool get isUsingFirebase => currentBackend == ApiBackend.firebase;
}

// Environment-specific configuration
class EnvironmentConfig {
  static bool get isDevelopment => F.appFlavor == Flavor.DEV;
  static bool get isProduction => F.appFlavor == Flavor.PROD;

  // Development vs Production API URLs for Go backend
  static String get goApiBaseUrl {
    if (isDevelopment) {
      return 'https://api-staging.chaseapp.tv';
    } else {
      return 'https://api.chaseapp.tv';
    }
  }
}