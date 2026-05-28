import 'package:flutter/material.dart';

class PrivacyScreen extends StatelessWidget {
  const PrivacyScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Privacy')),
      body: const ListView(
        padding: EdgeInsets.all(16),
        children: [
          ListTile(leading: Icon(Icons.visibility_off_outlined), title: Text('Anonymous sightings'), subtitle: Text('Owners see coarse areas, not finder identity.')),
          ListTile(leading: Icon(Icons.storefront_outlined), title: Text('Merchant zones'), subtitle: Text('Exact merchant identity is hidden unless recovery assistance is accepted.')),
          ListTile(leading: Icon(Icons.delete_outline), title: Text('Data controls'), subtitle: Text('Disable finder participation, delete tags, or delete your account from settings.')),
        ],
      ),
    );
  }
}
