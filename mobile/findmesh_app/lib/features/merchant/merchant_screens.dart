import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/app_state.dart';
import '../stand_provisioning/stand_provisioning_screen.dart';

class MerchantHomeScreen extends ConsumerWidget {
  const MerchantHomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        const ListTile(
          leading: Icon(Icons.storefront_outlined),
          title: Text('Merchant mode'),
          subtitle: Text('Manage counter stands and anonymous recovery requests.'),
        ),
        FilledButton.icon(
          icon: const Icon(Icons.add_business_outlined),
          label: const Text('Create merchant profile'),
          onPressed: () async {
            await ref.read(apiClientProvider).post('/v1/merchants', {
              'name': 'Demo Merchant',
              'display_name': 'Demo Merchant',
              'city': 'Bengaluru',
              'category': 'retail',
              'display_area': 'near a participating merchant zone',
            });
          },
        ),
        const SizedBox(height: 8),
        OutlinedButton.icon(
          icon: const Icon(Icons.nfc_outlined),
          label: const Text('Claim stand'),
          onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => const ClaimStandScreen())),
        ),
      ],
    );
  }
}

class StandHealthScreen extends StatelessWidget {
  const StandHealthScreen({required this.standId, super.key});

  final String standId;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Stand health')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          ListTile(title: const Text('Stand'), subtitle: Text(standId)),
          const ListTile(title: Text('Health model'), subtitle: Text('Heartbeat, battery, Wi-Fi RSSI, buffer count, firmware, and last error.')),
        ],
      ),
    );
  }
}
