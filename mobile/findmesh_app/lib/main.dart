import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'core/app_state.dart';
import 'features/admin_debug/admin_debug_screen.dart';
import 'features/anti_stalking/anti_stalking_screen.dart';
import 'features/auth/auth_screen.dart';
import 'features/merchant/merchant_screens.dart';
import 'features/privacy/privacy_screen.dart';
import 'features/recovery/recovery_screen.dart';
import 'features/settings/settings_screen.dart';
import 'features/tags/tag_screens.dart';

void main() {
  runApp(const ProviderScope(child: FindMeshApp()));
}

class FindMeshApp extends ConsumerWidget {
  const FindMeshApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final session = ref.watch(sessionProvider);
    return MaterialApp(
      title: 'FindMesh',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF16695C)),
        useMaterial3: true,
        visualDensity: VisualDensity.standard,
      ),
      home: session.token == null ? const LoginScreen() : const HomeScreen(),
    );
  }
}

class HomeScreen extends ConsumerStatefulWidget {
  const HomeScreen({super.key});

  @override
  ConsumerState<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends ConsumerState<HomeScreen> {
  int index = 0;

  @override
  Widget build(BuildContext context) {
    final pages = [
      const MyTagsScreen(),
      const MerchantHomeScreen(),
      const RecoveryRequestsScreen(),
      const AntiStalkingAlertsScreen(),
      const SettingsScreen(),
    ];
    return Scaffold(
      appBar: AppBar(
        title: const Text('FindMesh'),
        actions: [
          IconButton(
            tooltip: 'Privacy',
            icon: const Icon(Icons.privacy_tip_outlined),
            onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => const PrivacyScreen())),
          ),
          IconButton(
            tooltip: 'Debug',
            icon: const Icon(Icons.bug_report_outlined),
            onPressed: () => Navigator.of(context).push(MaterialPageRoute(builder: (_) => const AdminDebugScreen())),
          ),
        ],
      ),
      body: pages[index],
      bottomNavigationBar: NavigationBar(
        selectedIndex: index,
        onDestinationSelected: (value) => setState(() => index = value),
        destinations: const [
          NavigationDestination(icon: Icon(Icons.sell_outlined), label: 'Tags'),
          NavigationDestination(icon: Icon(Icons.storefront_outlined), label: 'Merchant'),
          NavigationDestination(icon: Icon(Icons.forum_outlined), label: 'Recovery'),
          NavigationDestination(icon: Icon(Icons.shield_outlined), label: 'Safety'),
          NavigationDestination(icon: Icon(Icons.settings_outlined), label: 'Settings'),
        ],
      ),
    );
  }
}
