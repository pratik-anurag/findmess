import 'dart:async';
import 'dart:typed_data';

import 'package:flutter_ble_peripheral/flutter_ble_peripheral.dart';
import 'package:flutter_blue_plus/flutter_blue_plus.dart';
import 'package:permission_handler/permission_handler.dart';

import 'ble_service.dart';
import 'findmesh_ble_protocol.dart';

class RealBleService implements BleService {
  RealBleService({FlutterBlePeripheral? peripheral}) : _peripheral = peripheral ?? FlutterBlePeripheral();

  final FlutterBlePeripheral _peripheral;

  @override
  Future<bool> isAvailable() async {
    final supported = await FlutterBluePlus.isSupported;
    final peripheralSupported = await _peripheral.isSupported;
    return supported || peripheralSupported;
  }

  @override
  Future<void> requestPermissions() async {
    await [
      Permission.bluetoothScan,
      Permission.bluetoothConnect,
      Permission.bluetoothAdvertise,
      Permission.locationWhenInUse,
    ].request();
    if (await FlutterBluePlus.isSupported) {
      try {
        await FlutterBluePlus.turnOn();
      } catch (_) {
        // iOS does not allow apps to turn Bluetooth on programmatically.
      }
    }
    await _peripheral.requestPermission();
  }

  @override
  Stream<BleAdvertisement> scanFindMeshTags() => _scanFor('FM_TAG');

  @override
  Stream<BleAdvertisement> scanMerchantZones() => _scanFor('FM_ZONE');

  @override
  Future<void> connectToTag(String ephemeralId) async {
    await requestPermissions();
    await for (final adv in scanFindMeshTags()) {
      if (adv.ephemeralId == ephemeralId && adv.deviceId != null) {
        await FlutterBluePlus.stopScan();
        return;
      }
    }
  }

  @override
  Future<void> connectToStandProvisioning(String standId) async {
    await requestPermissions();
  }

  @override
  Future<void> ringNearbyTag(String tagId) async {
    await requestPermissions();
  }

  @override
  Future<void> startHackathonTagAdvertiser(String ephemeralId) async {
    final normalized = FindMeshBleProtocol.normalizeEphemeralId(ephemeralId);
    await _startAdvertiser(
      localName: 'FM_TAG:${normalized.substring(0, 8)}',
      manufacturerData: FindMeshBleProtocol.encodeTag(normalized, lostHint: true),
    );
  }

  @override
  Future<void> startHackathonZoneAdvertiser(String zoneEphemeralId) async {
    final normalized = FindMeshBleProtocol.normalizeEphemeralId(zoneEphemeralId);
    await _startAdvertiser(
      localName: 'FM_ZONE:${normalized.substring(0, 8)}',
      manufacturerData: FindMeshBleProtocol.encodeZone(normalized),
    );
  }

  @override
  Future<void> stopAdvertising() => _peripheral.stop();

  Stream<BleAdvertisement> _scanFor(String type) async* {
    await requestPermissions();
    if (!await FlutterBluePlus.isSupported) return;
    await FlutterBluePlus.adapterState.where((state) => state == BluetoothAdapterState.on).first.timeout(const Duration(seconds: 8));
    await FlutterBluePlus.stopScan();
    await FlutterBluePlus.startScan(timeout: const Duration(seconds: 20), removeIfGone: const Duration(seconds: 5));
    try {
      final emitted = <String>{};
      await for (final results in FlutterBluePlus.onScanResults) {
        for (final result in results) {
          final parsed = FindMeshBleProtocol.parse(
            serviceData: result.advertisementData.serviceData.values,
            manufacturerData: result.advertisementData.manufacturerData.entries
                .where((entry) => entry.key == FindMeshBleProtocol.manufacturerId)
                .map((entry) => entry.value),
            rssi: result.rssi,
            deviceId: result.device.remoteId.toString(),
            name: result.advertisementData.advName,
          );
          if (parsed == null || parsed.type != type) continue;
          final key = '${parsed.type}:${parsed.ephemeralId}:${parsed.deviceId}:${parsed.rssi}';
          if (emitted.add(key)) yield parsed;
        }
      }
    } finally {
      await FlutterBluePlus.stopScan();
    }
  }

  Future<void> _startAdvertiser({required String localName, required List<int> manufacturerData}) async {
    await requestPermissions();
    final state = await _peripheral.start(
      advertiseData: AdvertiseData(
        manufacturerId: FindMeshBleProtocol.manufacturerId,
        manufacturerData: Uint8List.fromList(manufacturerData),
        localName: localName,
        includeDeviceName: false,
      ),
      advertiseSettings: AdvertiseSettings(
        connectable: false,
        timeout: 1800,
        advertiseMode: AdvertiseMode.advertiseModeLowLatency,
        txPowerLevel: AdvertiseTxPower.advertiseTxPowerMedium,
      ),
    );
    if (state.toString().toLowerCase().contains('unsupported')) {
      throw StateError('BLE advertising is not supported on this device');
    }
  }
}
