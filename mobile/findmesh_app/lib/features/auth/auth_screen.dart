import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/app_state.dart';

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  final phone = TextEditingController();
  final otp = TextEditingController();
  bool otpSent = false;
  String? error;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('FindMesh')),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          const Text('Recover lost items privately using anonymous sightings from your devices and participating merchant stands.'),
          const SizedBox(height: 24),
          TextField(controller: phone, keyboardType: TextInputType.phone, decoration: const InputDecoration(labelText: 'Phone number')),
          if (otpSent) ...[
            const SizedBox(height: 12),
            TextField(controller: otp, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: 'OTP')),
          ],
          if (error != null) Padding(padding: const EdgeInsets.only(top: 12), child: Text(error!, style: TextStyle(color: Theme.of(context).colorScheme.error))),
          const SizedBox(height: 20),
          FilledButton.icon(
            icon: Icon(otpSent ? Icons.login : Icons.sms_outlined),
            label: Text(otpSent ? 'Verify OTP' : 'Send OTP'),
            onPressed: () async {
              try {
                if (!otpSent) {
                  await ref.read(sessionProvider.notifier).startOtp(phone.text);
                  setState(() => otpSent = true);
                } else {
                  await ref.read(sessionProvider.notifier).verifyOtp(phone.text, otp.text);
                }
              } catch (e) {
                setState(() => error = e.toString());
              }
            },
          ),
          if (enableDevLogin) ...[
            const SizedBox(height: 12),
            OutlinedButton.icon(
              icon: const Icon(Icons.bolt_outlined),
              label: const Text('Continue with demo login'),
              onPressed: () async {
                try {
                  await ref.read(sessionProvider.notifier).devLogin();
                } catch (e) {
                  setState(() => error = e.toString());
                }
              },
            ),
          ],
        ],
      ),
    );
  }
}
