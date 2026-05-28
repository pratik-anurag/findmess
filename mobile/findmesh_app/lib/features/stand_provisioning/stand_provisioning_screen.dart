import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/app_state.dart';

class ClaimStandScreen extends ConsumerStatefulWidget {
  const ClaimStandScreen({super.key});

  @override
  ConsumerState<ClaimStandScreen> createState() => _ClaimStandScreenState();
}

class _ClaimStandScreenState extends ConsumerState<ClaimStandScreen> {
  final serial = TextEditingController(text: 'FM-STAND-DEV-1');
  String? standId;
  String? token;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Claim stand')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          const Text('Tap NFC or scan the setup QR, then connect to BLE provisioning to configure Wi-Fi.'),
          TextField(controller: serial, decoration: const InputDecoration(labelText: 'Stand serial')),
          const SizedBox(height: 12),
          FilledButton.icon(
            icon: const Icon(Icons.nfc_outlined),
            label: const Text('Start claim'),
            onPressed: () async {
              final response = await ref.read(apiClientProvider).post('/v1/stands/claim/start', {
                'serial': serial.text,
                'public_key': '',
              });
              setState(() {
                standId = (response['stand'] as Map<String, dynamic>)['id'] as String;
                token = response['claim_token'] as String;
              });
            },
          ),
          if (standId != null) ...[
            ListTile(title: const Text('Stand ID'), subtitle: Text(standId!)),
            FilledButton.icon(
              icon: const Icon(Icons.wifi_outlined),
              label: const Text('Configure Wi-Fi'),
              onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => StandWifiProvisioningScreen(standId: standId!, claimToken: token!))),
            ),
          ],
        ],
      ),
    );
  }
}

class StandWifiProvisioningScreen extends ConsumerStatefulWidget {
  const StandWifiProvisioningScreen({required this.standId, required this.claimToken, super.key});

  final String standId;
  final String claimToken;

  @override
  ConsumerState<StandWifiProvisioningScreen> createState() => _StandWifiProvisioningScreenState();
}

class _StandWifiProvisioningScreenState extends ConsumerState<StandWifiProvisioningScreen> {
  final ssid = TextEditingController();
  final password = TextEditingController();
  final merchantId = TextEditingController(text: '00000000-0000-4000-8000-000000000201');
  final zoneId = TextEditingController(text: '00000000-0000-4000-8000-000000000202');

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Stand Wi-Fi')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          TextField(controller: ssid, decoration: const InputDecoration(labelText: 'Wi-Fi SSID')),
          TextField(controller: password, obscureText: true, decoration: const InputDecoration(labelText: 'Wi-Fi password')),
          TextField(controller: merchantId, decoration: const InputDecoration(labelText: 'Merchant ID')),
          TextField(controller: zoneId, decoration: const InputDecoration(labelText: 'Zone ID')),
          const SizedBox(height: 16),
          FilledButton.icon(
            icon: const Icon(Icons.check_circle_outline),
            label: const Text('Finish provisioning'),
            onPressed: () async {
              await ref.read(bleServiceProvider).connectToStandProvisioning(widget.standId);
              await ref.read(apiClientProvider).post('/v1/stands/claim/complete', {
                'stand_id': widget.standId,
                'token': widget.claimToken,
                'merchant_id': merchantId.text,
                'zone_id': zoneId.text,
              });
              if (context.mounted) Navigator.of(context).popUntil((route) => route.isFirst);
            },
          ),
        ],
      ),
    );
  }
}
