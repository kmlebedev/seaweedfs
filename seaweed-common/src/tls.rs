//! The process-wide rustls crypto provider.
//!
//! Both binaries link more than one provider — in the volume server aws-lc-rs
//! comes from rustls' own defaults and ring from reqwest, in the lance worker
//! from lance's `aws` backend and reqwest — so rustls cannot auto-select one,
//! and tonic's client TLS panics on first use until something is installed.
//! Each binary therefore pins one, and it has to be the same one, which is why
//! the choice lives here rather than in either tree.
//!
//! Which provider that is became a build choice with the e2k port (issue
//! #13007): aws-lc-rs has no e2k backend, so the e2k build selects
//! `crypto-ring`. The two features are mutually exclusive, and one of them has
//! to be on.

use rustls::crypto::CryptoProvider;

#[cfg(all(feature = "crypto-aws-lc", feature = "crypto-ring"))]
compile_error!(
    "seaweed-common: `crypto-aws-lc` and `crypto-ring` cannot both be enabled; \
     pick the provider the target can build"
);

#[cfg(not(any(feature = "crypto-aws-lc", feature = "crypto-ring")))]
compile_error!(
    "seaweed-common: enable exactly one of `crypto-aws-lc` or `crypto-ring`"
);

/// The provider this build was configured with.
pub fn default_provider() -> CryptoProvider {
    #[cfg(feature = "crypto-aws-lc")]
    let provider = rustls::crypto::aws_lc_rs::default_provider();
    #[cfg(feature = "crypto-ring")]
    let provider = rustls::crypto::ring::default_provider();
    provider
}

/// Pin rustls's process-wide default provider to `default_provider()`.
/// Idempotent: the first call wins and every later one is a no-op, so callers
/// do not have to coordinate.
pub fn install_default_crypto_provider() {
    let _ = default_provider().install_default();
}

#[cfg(test)]
mod tests {
    use super::install_default_crypto_provider;

    #[test]
    fn installing_is_idempotent_and_leaves_a_default_behind() {
        install_default_crypto_provider();
        install_default_crypto_provider();
        assert!(rustls::crypto::CryptoProvider::get_default().is_some());
    }
}
