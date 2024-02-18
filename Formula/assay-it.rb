# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class AssayIt < Formula
  desc "Confirm Quality and Eliminate Risk by Testing Microservices in Production."
  homepage "https://assay.it"
  version "1.2.10"
  license "MIT"

  depends_on "go"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/assay-it/assay-it/releases/download/v1.2.10/assay-it_1.2.10_darwin_amd64"
      sha256 "3af6d65893dcb745d34789cbf3ac371528d135bcb55d95b7bd96ee02dceebdab"

      def install
        bin.install "assay-it_1.2.10_darwin_amd64" => "assay-it"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/assay-it/assay-it/releases/download/v1.2.10/assay-it_1.2.10_darwin_arm64"
      sha256 "8911e7d6268acdd9196820d778e3e514687108a80a5e8f54e71587628031abb5"

      def install
        bin.install "assay-it_1.2.10_darwin_arm64" => "assay-it"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/assay-it/assay-it/releases/download/v1.2.10/assay-it_1.2.10_linux_amd64"
      sha256 "337ba6cc24d7299feb095d3cca6d8fb6d103d135c5b76e238e98ffeaf8198454"

      def install
        bin.install "assay-it_1.2.10_linux_amd64" => "assay-it"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/assay-it/assay-it/releases/download/v1.2.10/assay-it_1.2.10_linux_arm64"
      sha256 "431fdf02f473cd35556b4d099dfdce43e518fc816ccba738af31ec730e7a96bf"

      def install
        bin.install "assay-it_1.2.10_linux_arm64" => "assay-it"
      end
    end
  end

  test do
    system "#{bin}/assay-it -v"
  end
end
