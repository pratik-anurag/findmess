import 'dart:convert';

import 'package:http/http.dart' as http;

class ApiException implements Exception {
  ApiException(this.message, this.statusCode);

  final String message;
  final int statusCode;

  @override
  String toString() => 'ApiException($statusCode): $message';
}

class ApiClient {
  ApiClient({
    this.baseUrl = const String.fromEnvironment('FINDMESH_API_BASE_URL', defaultValue: 'http://localhost:8080'),
    http.Client? httpClient,
  }) : _http = httpClient ?? http.Client();

  final String baseUrl;
  final http.Client _http;
  String? token;

  Future<Map<String, dynamic>> post(String path, Map<String, dynamic> body) async {
    return _send('POST', path, body: body);
  }

  Future<Map<String, dynamic>> patch(String path, Map<String, dynamic> body) async {
    return _send('PATCH', path, body: body);
  }

  Future<dynamic> get(String path) => _send('GET', path);

  Future<dynamic> delete(String path) => _send('DELETE', path);

  Future<dynamic> _send(String method, String path, {Map<String, dynamic>? body}) async {
    final uri = Uri.parse('$baseUrl$path');
    final headers = <String, String>{'content-type': 'application/json'};
    if (token != null) {
      headers['authorization'] = 'Bearer $token';
    }
    for (var attempt = 0; attempt < 3; attempt++) {
      final response = await _request(method, uri, headers, body);
      if (response.statusCode >= 500 && attempt < 2) {
        await Future<void>.delayed(Duration(milliseconds: 150 * (attempt + 1)));
        continue;
      }
      final decoded = response.body.isEmpty ? <String, dynamic>{} : jsonDecode(response.body);
      if (response.statusCode >= 400) {
        final message = decoded is Map<String, dynamic> ? decoded['error'] as String? : null;
        throw ApiException(message ?? 'Request failed', response.statusCode);
      }
      return decoded;
    }
    throw ApiException('Request failed after retries', 0);
  }

  Future<http.Response> _request(String method, Uri uri, Map<String, String> headers, Map<String, dynamic>? body) {
    final raw = body == null ? null : jsonEncode(body);
    switch (method) {
      case 'GET':
        return _http.get(uri, headers: headers);
      case 'PATCH':
        return _http.patch(uri, headers: headers, body: raw);
      case 'DELETE':
        return _http.delete(uri, headers: headers);
      default:
        return _http.post(uri, headers: headers, body: raw);
    }
  }
}
