# Philosophy

Stratus Red Team is opinionated about the attack techniques it packages, in order to make sure it provides actual value, as opposed to emulating "attacker behavior" in a non-actionable way (such as `calling sts:GetCallerIdentity`).

This page describes the characteristics that all attack techniques of Stratus Red Team should have.

## Be Granular

An attack technique should be **granular**, meaning that it should emulate a single step of an attack.

- Good: Share an EBS snapshot with an external AWS account.
- Bad: Use an IAM access key to perform privilege escalation, run discovery commands, take an EBS snapshot of an instance, share the EBS snapshot with an external AWS account.

## Emulate plausible attacker activity

It's always hard to draw a line between legitimate and malicious activity. Attack techniques should emulate **plausible attacker activity**, i.e. activity that can legitimately be thought to be malicious.

It should always be possible to derive a reasonable detection rule from a Stratus Red Team attack technique.

- Good: Delete a CloudTrail trail
- Bad: Run `sts:GetCallerIdentity`
    - While attackers might use this API call, it is in no way indicative of attacker activity, as it's used by many services and client applications.
    - It isn't emulating activity that could reasonably be thought to be malicious.

Stratus Red Team's goal is *not* to re-implement all AWS API calls that may be used by an attacker.

## Be Self-Sufficient

An attack technique should not be dependent on the state of the cloud environment it's run against.

- Good: Create an EBS snapshot and share it
- Bad: Expect an EBS snapshot exists in the account prior to running Stratus Red Team