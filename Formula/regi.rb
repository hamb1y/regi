class Regi < Formula
  desc "Tiny newline-delimited plaintext register CLI"
  homepage "https://github.com/hamb1y/regi"
  url "https://github.com/hamb1y/regi/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "SHA256_FOR_V0_1_0_TARBALL"
  head "https://github.com/hamb1y/regi.git", branch: "main"
  license "BSD-3-Clause"

  depends_on "go" => :build

  def install
    system "go", "build", "-trimpath", "-ldflags=-s -w", "-o", bin/"regi", "."
  end

  test do
    ENV["HOME"] = testpath
    system bin/"regi", "add", "test", "hello"
    assert_equal "hello\n", shell_output("#{bin}/regi test")
    assert_path_exists testpath/".config/regi/registers/test.regi"
  end
end
