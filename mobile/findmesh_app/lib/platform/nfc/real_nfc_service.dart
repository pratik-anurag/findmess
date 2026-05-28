import 'dart:async';
import 'dart:convert';

import 'package:nfc_manager/nfc_manager.dart';

import 'nfc_service.dart';

class RealNfcService implements NfcService {
  @override
  Future<bool> isAvailable() => NfcManager.instance.isAvailable();

  @override
  Future<String?> readPayload() async {
    if (!await isAvailable()) return null;
    final completer = Completer<String?>();
    await NfcManager.instance.startSession(
      onDiscovered: (tag) async {
        try {
          final ndef = Ndef.from(tag);
          final message = ndef?.cachedMessage;
          final payload = message == null ? null : _decodeFirstRecord(message.records);
          await NfcManager.instance.stopSession();
          if (!completer.isCompleted) completer.complete(payload);
        } catch (error) {
          await NfcManager.instance.stopSession(errorMessage: error.toString());
          if (!completer.isCompleted) completer.completeError(error);
        }
      },
    );
    return completer.future.timeout(const Duration(seconds: 25), onTimeout: () async {
      await NfcManager.instance.stopSession(errorMessage: 'Timed out waiting for NFC tag');
      return null;
    });
  }

  @override
  Future<void> writePayload(String payload) async {
    if (!await isAvailable()) {
      throw StateError('NFC is not available on this device');
    }
    final completer = Completer<void>();
    await NfcManager.instance.startSession(
      onDiscovered: (tag) async {
        try {
          final ndef = Ndef.from(tag);
          if (ndef == null) {
            throw StateError('NFC tag is not NDEF compatible');
          }
          if (!ndef.isWritable) {
            throw StateError('NFC tag is read only');
          }
          final uri = Uri.tryParse(payload);
          final record = uri == null ? NdefRecord.createText(payload) : NdefRecord.createUri(uri);
          await ndef.write(NdefMessage([record]));
          await NfcManager.instance.stopSession();
          if (!completer.isCompleted) completer.complete();
        } catch (error) {
          await NfcManager.instance.stopSession(errorMessage: error.toString());
          if (!completer.isCompleted) completer.completeError(error);
        }
      },
    );
    return completer.future.timeout(const Duration(seconds: 25), onTimeout: () async {
      await NfcManager.instance.stopSession(errorMessage: 'Timed out waiting for NFC tag');
      throw TimeoutException('Timed out waiting for NFC tag');
    });
  }

  String? _decodeFirstRecord(List<NdefRecord> records) {
    for (final record in records) {
      final type = ascii.decode(record.type, allowInvalid: true);
      if (type == 'U' && record.payload.isNotEmpty) {
        return '${_uriPrefix(record.payload.first)}${utf8.decode(record.payload.sublist(1), allowMalformed: true)}';
      }
      if (type == 'T' && record.payload.isNotEmpty) {
        final languageLength = record.payload.first & 0x3f;
        final textStart = 1 + languageLength;
        if (record.payload.length >= textStart) {
          return utf8.decode(record.payload.sublist(textStart), allowMalformed: true);
        }
      }
      final raw = utf8.decode(record.payload, allowMalformed: true);
      if (raw.contains('findmesh://')) {
        return raw.substring(raw.indexOf('findmesh://'));
      }
    }
    return null;
  }

  String _uriPrefix(int code) {
    const prefixes = <int, String>{
      0x00: '',
      0x01: 'http://www.',
      0x02: 'https://www.',
      0x03: 'http://',
      0x04: 'https://',
    };
    return prefixes[code] ?? '';
  }
}
