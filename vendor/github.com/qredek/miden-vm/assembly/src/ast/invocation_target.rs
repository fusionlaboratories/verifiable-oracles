use super::{parsers::decode_hex_rpo_digest_label, LibraryPath, ParsingError, RpoDigest, Token};

/// Describes targets of `exec`, `call`, and `syscall` instructions.
pub enum InvocationTarget<'a> {
    MastRoot(RpoDigest),
    ProcedureName(&'a str),
    ProcedurePath { name: &'a str, module: &'a str },
}

impl<'a> InvocationTarget<'a> {
    /// Parses the provided label into an invocation target.
    ///
    /// A label of an invoked procedure must comply with the following rules:
    /// - It can be a hexadecimal string representing a MAST root digest ([RpoDigest]). In this case,
    ///   the label must start with "0x" and must be followed by a valid hexadecimal string
    ///   representation of an [RpoDigest].
    /// - It can contain a single procedure name. In this case, the label must comply with procedure
    ///   name rules.
    /// - It can contain module name followed by procedure name (e.g., "module::procedure"). In this
    ///   case both module and procedure name must comply with relevant name rules.
    ///
    /// All other combinations will result in an error.
    pub fn parse(label: &'a str, token: &'a Token) -> Result<Self, ParsingError> {
        if label.starts_with("0x") {
            return Ok(InvocationTarget::MastRoot(
                decode_hex_rpo_digest_label(label)
                    .map_err(|err| ParsingError::invalid_proc_root_invocation(token, label, err))?,
            ));
        }

        let num_components = LibraryPath::validate(label)
            .map_err(|_| ParsingError::invalid_proc_invocation(token, label))?;

        match num_components {
            1 => Ok(InvocationTarget::ProcedureName(label)),
            2 => {
                let parts = label.split_once(LibraryPath::PATH_DELIM).expect("no components");
                Ok(InvocationTarget::ProcedurePath {
                    name: parts.1,
                    module: parts.0,
                })
            }
            _ => Err(ParsingError::invalid_proc_invocation(token, label)),
        }
    }
}
