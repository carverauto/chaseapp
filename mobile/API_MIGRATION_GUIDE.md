# Mobile Client API Migration Guide

This guide outlines the steps to migrate the Flutter mobile app from Firebase backend to the new Go monolithic API.

## Overview

The mobile app currently uses Firebase SDK for most operations (Firestore, Authentication, Cloud Functions). The migration will transition to HTTP API calls while maintaining Firebase Authentication for user management.

## Migration Approach

### Phase 1: Configuration Setup
1. **API Configuration**: Use the new `ApiConfig` class in `lib/src/const/api_config.dart`
2. **Backend Toggle**: Switch between Firebase and Go API using the `currentBackend` setting
3. **Environment Setup**: Configure staging/production API endpoints

### Phase 2: HTTP Client Implementation
1. **Add HTTP Client**: Install and configure Dio or HTTP package
2. **Create API Service**: Implement service classes for each API endpoint
3. **Handle Authentication**: Configure JWT token management

### Phase 3: Gradual Migration
1. **Start with Non-Critical Features**: Begin with weather, external data endpoints
2. **Core Features**: Migrate chases, aircraft, and other main features
3. **Real-time Features**: Implement WebSocket/NATS client for real-time updates

## Implementation Steps

### 1. Dependencies
Add to `pubspec.yaml`:
```yaml
dependencies:
  dio: ^5.3.0  # HTTP client
  json_annotation: ^4.8.0  # JSON serialization

dev_dependencies:
  json_serializable: ^6.7.0  # Code generation for JSON
  build_runner: ^2.4.6  # Code generation runner
```

### 2. API Service Structure
Create service classes for each API domain:

#### `lib/src/services/api_service.dart`
```dart
import 'package:dio/dio.dart';
import '../const/api_config.dart';

class ApiService {
  late final Dio _dio;

  ApiService() {
    _dio = Dio(BaseOptions(
      baseUrl: ApiConfig.baseUrl,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ));

    // Add authentication interceptor
    _dio.interceptors.add(AuthInterceptor());
  }

  Future<Response<T>> get<T>(String path, {Map<String, dynamic>? queryParameters}) {
    return _dio.get<T>(path, queryParameters: queryParameters);
  }

  Future<Response<T>> post<T>(String path, {dynamic data}) {
    return _dio.post<T>(path, data: data);
  }

  // Other HTTP methods...
}
```

### 3. Domain-Specific Services

#### Chases Service
`lib/src/services/chases_service.dart`:
```dart
class ChasesService {
  final ApiService _apiService = ApiService();

  Future<List<Chase>> getChases({int limit = 20, int offset = 0}) async {
    final response = await _apiService.get(
      ApiConfig.chasesEndpoint,
      queryParameters: {'limit': limit, 'offset': offset},
    );
    return (response.data as List).map((json) => Chase.fromJson(json)).toList();
  }

  Future<Chase> getChase(String id) async {
    final response = await _apiService.get('${ApiConfig.chasesEndpoint}/$id');
    return Chase.fromJson(response.data);
  }

  // Other chase operations...
}
```

#### Push Notifications Service
`lib/src/services/push_service.dart`:
```dart
class PushService {
  final ApiService _apiService = ApiService();

  Future<void> subscribeToTopic(String token, List<String> topics) async {
    await _apiService.post(ApiConfig.subscribePushEndpoint, data: {
      'token': token,
      'platform': Platform.isIOS ? 'ios' : 'android',
      'topics': topics,
    });
  }

  Future<void> unsubscribeFromTopic(String token, List<String> topics) async {
    await _apiService.post(ApiConfig.unsubscribePushEndpoint, data: {
      'token': token,
      'topics': topics,
    });
  }
}
```

### 4. Authentication Integration

#### Auth Interceptor
`lib/src/interceptors/auth_interceptor.dart`:
```dart
class AuthInterceptor extends Interceptor {
  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    // Add Firebase ID token to requests
    final user = FirebaseAuth.instance.currentUser;
    if (user != null) {
      user.getIdToken().then((token) {
        options.headers['Authorization'] = 'Bearer $token';
        handler.next(options);
      });
    } else {
      handler.next(options);
    }
  }
}
```

### 5. Model Updates

Update model classes to work with both Firestore and HTTP API:

#### Chase Model Example
```dart
import 'package:json_annotation/json_annotation.dart';

part 'chase.g.dart';

@JsonSerializable()
class Chase {
  final String id;
  final String title;
  final String? description;
  @JsonKey(name: 'chase_type')
  final String chaseType;
  final bool live;
  @JsonKey(name: 'started_at')
  final DateTime? startedAt;
  @JsonKey(name: 'ended_at')
  final DateTime? endedAt;
  // ... other fields

  Chase({
    required this.id,
    required this.title,
    this.description,
    required this.chaseType,
    required this.live,
    this.startedAt,
    this.endedAt,
    // ... other parameters
  });

  factory Chase.fromJson(Map<String, dynamic> json) => _$ChaseFromJson(json);
  factory Chase.fromFirestore(DocumentSnapshot doc) {
    final data = doc.data() as Map<String, dynamic>;
    return Chase.fromJson({...data, 'id': doc.id});
  }

  Map<String, dynamic> toJson() => _$ChaseToJson(this);
}
```

### 6. Repository Pattern Updates

Update repositories to use the appropriate backend:

#### Chase Repository
`lib/src/modules/chase/domain/chase_repo.dart`:
```dart
class ChaseRepository {
  final ChasesService? _chasesService;
  final CollectionReference<Chase> _firestoreCollection;

  ChaseRepository({ChasesService? chasesService})
      : _chasesService = chasesService,
        _firestoreCollection = chasesCollectionRef;

  Future<List<Chase>> getChases({int limit = 20, int offset = 0}) async {
    if (ApiConfig.isUsingGoApi && _chasesService != null) {
      return await _chasesService!.getChases(limit: limit, offset: offset);
    } else {
      final snapshot = await _firestoreCollection
          .orderBy('created_at', descending: true)
          .limit(limit)
          .get();
      return snapshot.docs.map((doc) => doc.data()).toList();
    }
  }

  // Other repository methods...
}
```

## Migration Checklist

### Phase 1: Preparation
- [ ] Add HTTP client dependencies
- [ ] Create API configuration
- [ ] Set up code generation for JSON serialization
- [ ] Create base API service class

### Phase 2: Service Implementation
- [ ] Implement authentication service
- [ ] Create chases service
- [ ] Create aircraft service
- [ ] Create weather service
- [ ] Create push notifications service
- [ ] Create external data services (quakes, boats, launches)

### Phase 3: Integration
- [ ] Update repositories to use HTTP services
- [ ] Update models for JSON serialization
- [ ] Add error handling for HTTP responses
- [ ] Implement retry logic for network failures
- [ ] Add network connectivity checking

### Phase 4: Testing
- [ ] Unit test API services
- [ ] Integration test with staging API
- [ ] Test authentication flow
- [ ] Test error scenarios
- [ ] Performance testing

### Phase 5: Deployment
- [ ] Switch to Go API in staging
- [ ] Monitor for issues
- [ ] Gradual rollout to production
- [ ] Monitor performance and errors
- [ ] Full cutover to Go API

## Configuration Toggle

To switch between backends, update `lib/src/const/api_config.dart`:

```dart
// Change this line to switch backends
static const ApiBackend _currentBackend = ApiBackend.goApi; // or ApiBackend.firebase
```

## Considerations

### Authentication
- Continue using Firebase Authentication for user management
- Pass Firebase ID token in Authorization header
- Implement token refresh logic

### Real-time Features
- Current: Firestore real-time listeners and Stream Chat
- New: WebSocket connections or periodic polling
- Consider implementing NATS client for real-time updates

### Offline Support
- Firestore provides built-in offline caching
- Implement local caching for HTTP responses
- Use SQLite or local storage for offline data

### Error Handling
- HTTP status codes vs Firebase error codes
- Network connectivity issues
- API rate limiting
- Authentication token expiration

### Performance
- HTTP requests may be faster than Firestore
- Implement request caching
- Use connection pooling
- Optimize bundle size by removing unused Firebase SDKs

## Testing Strategy

### Unit Tests
- Test API service methods
- Test JSON serialization/deserialization
- Test error handling

### Integration Tests
- Test against staging API
- Test authentication flow
- Test real-time features

### End-to-End Tests
- Test complete user flows
- Test offline scenarios
- Test error recovery

## Rollback Plan

If issues arise during migration:
1. Toggle `ApiConfig.currentBackend` back to `ApiBackend.firebase`
2. Re-deploy app with Firebase configuration
3. Monitor for issues
4. Investigate and fix Go API issues
5. Retry migration when ready