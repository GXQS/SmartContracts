# GXVM Architecture

GXVM executes bytecode deterministically with bounded memory, bounded call depth, explicit gas accounting, authenticated trie-backed state roots, snapshot/revert semantics, and deterministic replay support.

Execution flow:
1. Verify bytecode and ABI constraints.
2. Open snapshot and run interpreter loop with step limits.
3. Charge opcode + memory expansion gas each step.
4. Commit snapshot on success, revert on explicit revert or fault.
5. Emit receipt with root hash, gas usage, and return data.
