import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:findmesh_app/main.dart';

void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  testWidgets('renders login screen', (tester) async {
    await tester.pumpWidget(const FindMeshApp());
    expect(find.text('FindMesh'), findsWidgets);
    expect(find.text('Send OTP'), findsOneWidget);
  });
}
