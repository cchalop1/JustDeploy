"""Exercise installer completion without installing packages or touching the host."""
import os
from pathlib import Path
import subprocess
import unittest


class InstallerOnboardingTest(unittest.TestCase):
    def run_completion(self, *, probe="http", active="yes", ssh=True):
        source = (Path(__file__).resolve().parents[1] / "install.sh").read_text()
        completion = source[source.index("# Wait for a real API response"):]
        mocks = r'''
set -euo pipefail
systemctl() { [ "$TEST_ACTIVE" = yes ]; }
curl() {
  case "$*" in
    *https://*) [ "$TEST_PROBE" = https ] ;;
    *http://*) [ "$TEST_PROBE" = http ] ;;
    *) return 1 ;;
  esac
}
sleep() { :; }
hostname() { echo '192.0.2.50'; }
'''
        environment = {**os.environ, "TEST_PROBE": probe, "TEST_ACTIVE": active,
                       "SSH_CONNECTION": "198.51.100.20 40000 203.0.113.10 22" if ssh else "",
                       "SUDO_USER": "vps-user"}
        return subprocess.run(["bash", "-c", mocks + completion], env=environment,
                              text=True, capture_output=True)

    def test_install_prints_tunnel_and_account_creation_instructions(self):
        result = self.run_completion()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("ssh -N -L 5915:127.0.0.1:5915 vps-user@203.0.113.10", result.stdout)
        self.assertIn("http://localhost:5915", result.stdout)
        self.assertIn("Create your admin account on first use", result.stdout)
        self.assertNotIn("setup-token", result.stdout)

    def test_address_fallback_and_existing_account_instructions(self):
        result = self.run_completion(ssh=False)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("vps-user@192.0.2.50", result.stdout)
        self.assertIn("sign in with your existing account", result.stdout)

    def test_direct_tls_uses_existing_https_address(self):
        result = self.run_completion(probe="https")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("direct TLS detected", result.stdout)
        self.assertNotIn("Open http://localhost", result.stdout)

    def test_unavailable_api_or_stopped_service_never_reports_success(self):
        for options in ({"probe": "unavailable"}, {"active": "no"}):
            with self.subTest(options=options):
                result = self.run_completion(**options)
                self.assertNotEqual(result.returncode, 0)
                self.assertIn("Installation has not been confirmed", result.stdout)
                self.assertNotIn("Next: open JustDeploy", result.stdout)


if __name__ == "__main__":
    unittest.main()
