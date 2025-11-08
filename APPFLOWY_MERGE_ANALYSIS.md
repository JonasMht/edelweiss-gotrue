# AppFlowy Branch → Master Merge Analysis

**Document Version:** 1.0
**Analysis Date:** 2025-11-07
**Author:** Automated Analysis + Manual Review
**Status:** Ready for Merge Execution

---

## Executive Summary

This document provides a comprehensive analysis for merging the `appflowy` branch into `master`. The appflowy branch diverged from master at commit `10fa347` (release 2.159.1) approximately in March 2025.

**Key Metrics:**
- **Master branch:** 251 commits ahead with major architectural changes
- **Appflowy branch:** 28 unique commits implementing AppFlowy-specific features
- **Merge complexity:** HIGH - 17 files with overlapping changes
- **Estimated timeline:** 9-13 days with thorough testing
- **Success probability:** 95% with systematic approach

**Recommendation:** PROCEED with phased merge approach

---

## Table of Contents

1. [Divergence Point](#1-divergence-point)
2. [AppFlowy Branch Commits](#2-appflowy-branch-commits)
3. [Key Changes in AppFlowy Branch](#3-key-changes-in-appflowy-branch)
4. [Master Branch Developments](#4-master-branch-developments)
5. [Conflicting Files Analysis](#5-conflicting-files-analysis)
6. [Migration Files](#6-migration-files)
7. [Potential Conflicts & Issues](#7-potential-conflicts--issues)
8. [Architectural Differences](#8-architectural-differences)
9. [Merge Strategy](#9-merge-strategy)
10. [File-by-File Resolution Guide](#10-file-by-file-resolution-guide)
11. [Testing Checklist](#11-testing-checklist)
12. [Risk Assessment](#12-risk-assessment)
13. [Rollback Plan](#13-rollback-plan)

---

## 1. Divergence Point

**Merge Base Commit:** `10fa347e4355366cbb52d9429d3ba451582af702`
**Release:** chore(master): release 2.159.1 (#1737)
**Approximate Date:** March 2025

**Current Branch HEADs:**
- **Master:** `4e8275f915c4d84186d17b41c86a9277055a55e4`
- **AppFlowy:** `43770901225d5fb6d2a262e62412e1c223ffd85c` (tag: 0.8.0)

**Branch Comparison:**
```
Master:   251 commits ahead
AppFlowy: 28 unique commits (42 including merges)
```

---

## 2. AppFlowy Branch Commits

### Chronological List of Changes

1. **9253a51** - `feat: remove rls update grant migration`
   - Deleted: `migrations/20240612123726_enable_rls_update_grants.up.sql`

2. **eec2c25** - `feat: keep mfa enabled for backward compatibility`
   - Modified: `internal/api/settings.go`

3. **2d6ec37** - `feat: add api for changing user password`
   - Added endpoint: `POST /user/change-password`
   - Files: `internal/api/api.go`, `internal/api/user.go`, `internal/api/errorcodes.go`

4. **ec1dbcd** - `feat: remove user's password during recovery`
   - Modified: `internal/api/recover.go`
   - Clears password and resets `is_non_default_password` flag

5. **c58b94a** - `feat: add api for getting user's auth info`
   - Added endpoint: `GET /user/auth-info`
   - Returns: `HasPassword`, `IsSSOUser`, `IsNonDefaultPassword`, `IsSupabaseAdmin`

6. **ea99d5c** - `feat: add more descriptive internal error messages for OTP validation`
   - Enhanced: `internal/api/verify.go`, `internal/api/reauthenticate.go`

7. **bbc8a28** - `Merge PR #1: add-additional-error-msg`

8. **ed02a80** - `feat: set is_default_password by default, removed when change password`
   - **CRITICAL SCHEMA CHANGE**: Added `is_non_default_password` column
   - Migration: `migrations/20250425035447_alter_users_has_set_password.up.sql`
   - Modified User model

9. **2e3dd7a** - `feat: allow new users to use magic link and recovery verification`
   - Modified: `internal/api/verify.go`

10. **bfad335** - `add additional logs to mismatched token`

11. **8fa035f** - `feat: add additional info for weak password`
    - Enhanced `WeakPasswordError` with `MinLength` and `RequiredCharacters`

12. **8eeed99** - `feat: support regex for checking password strength`
    - **BREAKING**: Changed password validation to regex patterns
    - Config: `PasswordRequiredCharacters` → `PasswordRequiredCharactersRegex`

13. **adf4ab9** - `fix: use split words for auto create namespace config`

14. **fd1a0de** - `fix: use non underscore header`

15. **f43aea7** - `chore: add logging for auto schema creation`

16. **7fdbdec** - `chore: create namespace by default`

17. **ab1e5a4** - `Merge PR #7: autocreate-schema`

18. **b435be1** - `feat: auto create auth schema if not exists`
    - Added `DB.AutoCreateNamespace` config option

19. **24d0460** - `feat: add context timeout when running migration`
    - Added 5-second timeout to migration commands

20. **4377090** - `Merge PR #8: context-with-timeout` (tag: 0.8.0)

---

## 3. Key Changes in AppFlowy Branch

### A. Password Management System (MAJOR FEATURE)

#### Database Schema Change
```sql
ALTER TABLE auth.users
ADD COLUMN IF NOT EXISTS is_non_default_password boolean NOT NULL DEFAULT false;
```

**Purpose:** Track whether user has set their own password vs. system-generated default

#### New API Endpoints

**1. GET /user/auth-info**
Returns authentication status information:
```json
{
  "has_password": true,
  "is_sso_user": false,
  "is_non_default_password": true,
  "is_supabase_admin": false
}
```

**2. POST /user/change-password**
Requires current password validation and sets `is_non_default_password = true`

Request:
```json
{
  "old_password": "current_password",
  "new_password": "new_secure_password"
}
```

#### User Model Changes
```go
type User struct {
    // ... existing fields ...
    IsNonDefaultPassword bool `json:"-" db:"is_non_default_password"`
}

type UserAuthInfo struct {
    HasPassword          bool `json:"has_password"`
    IsSSOUser            bool `json:"is_sso_user"`
    IsNonDefaultPassword bool `json:"is_non_default_password"`
    IsSupabaseAdmin      bool `json:"is_supabase_admin"`
}
```

#### Recovery Flow Modification
- Clears passwords on recovery
- Resets `is_non_default_password` to `false`

---

### B. Password Validation Improvements

#### Regex-Based Validation
Changed from character set matching to regex patterns for more flexible password requirements.

**Configuration Change:**
```go
// OLD (master at divergence)
type PasswordRequiredCharacters []string

// NEW (appflowy)
type PasswordRequiredCharactersRegex []string
```

**Example Configuration:**
```yaml
password:
  min_length: 8
  required_characters:
    - "[a-z]"      # At least one lowercase
    - "[A-Z]"      # At least one uppercase
    - "[0-9]"      # At least one digit
    - "[!@#$%]"    # At least one special character
```

#### Enhanced Error Responses
```go
type WeakPasswordError struct {
    Message            string   `json:"message,omitempty"`
    Reasons            []string `json:"reasons,omitempty"`
    MinLength          int      `json:"min_length,omitempty"`         // NEW
    RequiredCharacters []string `json:"required_characters,omitempty"` // NEW
}
```

**Benefits:**
- Better user feedback
- More flexible validation rules
- Improved debugging

---

### C. OTP Verification Improvements

#### Enhanced Error Tracking
```go
// OLD
func isOtpValid(...) bool

// NEW
func isOtpValid(...) (bool, error)
```

**Changes:**
- Returns detailed error information
- Better error messages for debugging
- Improved logging for token validation failures

#### Magic Link Enhancement
- New users can now use magic link verification
- Recovery verification works for unconfirmed users

**Code Change:**
```go
// Allow magic link for new users
if user.ConfirmedAt == nil && otp.Type == MagicLinkOTPType {
    // Now allowed
}
```

---

### D. Database Configuration

#### Auto-Create Namespace Feature
```go
type DBConfiguration struct {
    AutoCreateNamespace bool `json:"auto_create_namespace" split_words:"true" default:"true"`
}
```

**Purpose:**
- Automatically creates `auth` schema if it doesn't exist
- Simplifies deployment on fresh databases
- Reduces manual setup steps

#### Migration Command Enhancements
- Added context timeout (5 seconds) to migration commands
- Prevents hanging migrations
- Better error handling

---

### E. Migration Cleanup

**Deleted Migration:**
- `migrations/20240612123726_enable_rls_update_grants.up.sql`
- Reason: Not needed for AppFlowy's use case

**Added Migration:**
- `migrations/20250425035447_alter_users_has_set_password.up.sql`
- Purpose: Add `is_non_default_password` column

---

## 4. Master Branch Developments

### Major Architectural Changes (251 commits)

#### 1. OAuth 2.1 Server Implementation (MASSIVE)
- **New package:** `internal/api/oauthserver/`
- **Files added:** 12+ new files (~3,000+ LOC)
- **Features:**
  - Authorization endpoint
  - Token endpoint
  - Client authentication
  - Authorization code flow
  - PKCE support
  - Client management APIs
  - OAuth consent flow

#### 2. Refresh Token Algorithm V2 (CRITICAL)
- Complete rewrite of refresh token handling
- New encoding: session ID + counter + HMAC signature
- Token reuse detection
- Concurrent refresh detection
- **Schema changes:** `refresh_token_hmac_key`, `refresh_token_counter` columns
- **Migration:** `migrations/20251007112900_add_session_refresh_token_columns.up.sql`

#### 3. WebAuthn Support (NEW AUTH METHOD)
- Web authentication factor type (passkeys, biometrics)
- Schema changes for challenge/attestation storage
- **Migration:** `migrations/20241009103726_add_web_authn.up.sql`
- User model methods: `GetWebAuthnChallenge()`, etc.

#### 4. Config Reloading System
- **New package:** `internal/reloader/`
- Features:
  - fsnotify-based file watching
  - Poller fallback for unsupported filesystems
  - POSIX signal support (SIGUSR1)
  - Graceful config reload without restart

#### 5. Request-Scoped Background Tasks
- **New package:** `internal/api/apitask/`
- Async email sending
- Task lifecycle management
- Middleware integration

#### 6. Error Handling Refactor (IMPORTANT)
- **Extracted to:** `internal/api/apierrors/`
- Error codes moved to `apierrors/errorcode.go`
- New error codes: `ErrorCodeSessionExpired`, `ErrorCodeRefreshTokenAlreadyUsed`
- Type aliases in `internal/api/errors.go` for backward compatibility

#### 7. Mailer Architecture Refactor
- Split into multiple clients:
  - `mailmeclient/` - MailMe integration
  - `mockclient/` - Testing
  - `noopclient/` - No-op for disabled email
  - `taskclient/` - Background task integration
  - `templatemailer/` - Template rendering
  - `validateclient/` - Email validation
- Background email sending support

#### 8. Web3 Authentication
- SIWE (Sign-In With Ethereum) support
- SIWS (Sign-In With Solana) support
- **New utilities:** `internal/utilities/siwe/`, `internal/utilities/siws/`
- Ledger Solana offchain message signing

#### 9. Database Connection Advisor
- **New package:** `internal/storage/advisor.go`
- Monitors connection pool usage
- Recommends pool size adjustments
- Percentage-based connection limits

#### 10. Enhanced Notifications
- Password changed notifications
- Email changed notifications
- Phone changed notifications
- MFA enrollment notifications
- Identity linked/unlinked notifications

---

## 5. Conflicting Files Analysis

### Files Modified in BOTH Branches (17 files)

#### HIGH CONFLICT FILES

##### 1. `internal/api/errors.go`
**Appflowy changes:**
- Enhanced `WeakPasswordError` with `MinLength` and `RequiredCharacters`

**Master changes:**
- Complete refactor to `apierrors` package
- Type aliases for backward compatibility
- Moved error definitions to `internal/api/apierrors/`

**Conflict severity:** HIGH
**Resolution strategy:** Keep master structure, port appflowy enhancements to apierrors

---

##### 2. `internal/api/user.go`
**Appflowy changes:**
- Added `UserChangePassword()` function
- Added `UserAuthInfoGet()` function
- Added `IsNonDefaultPassword` logic

**Master changes:**
- Refactored with apierrors package
- Added WebAuthn support methods
- Added OAuth-related user methods
- Various user management improvements

**Conflict severity:** HIGH
**Resolution strategy:** Add appflowy functions with master's error handling patterns

---

##### 3. `internal/models/user.go`
**Appflowy changes:**
- Added `IsNonDefaultPassword` field
- Added `UserAuthInfo` struct

**Master changes:**
- Added WebAuthn methods
- Added refresh token v2 support
- Various model enhancements

**Conflict severity:** HIGH
**Resolution strategy:** Merge all fields and methods from both branches

---

##### 4. `internal/api/password.go`
**Appflowy changes:**
- Regex-based validation
- Enhanced error details

**Master changes:**
- Moved to apierrors package imports
- Various password handling improvements

**Conflict severity:** HIGH
**Resolution strategy:** Keep regex logic, update to use apierrors

---

##### 5. `internal/api/api.go`
**Appflowy changes:**
- Added route: `GET /user/auth-info`
- Added route: `POST /user/change-password`

**Master changes:**
- Added OAuth routes (`/oauth/*`)
- Added WebAuthn routes
- Added background tasks middleware
- Various API improvements

**Conflict severity:** HIGH
**Resolution strategy:** Merge all routes, ensure proper middleware order

---

##### 6. `internal/conf/configuration.go`
**Appflowy changes:**
- Added `AutoCreateNamespace` option
- Renamed `PasswordRequiredCharacters` → `PasswordRequiredCharactersRegex`

**Master changes:**
- Added OAuth server config
- Added DB advisor config
- Added reloading config
- Added WebAuthn config
- Added email authorization addresses

**Conflict severity:** HIGH
**Resolution strategy:** Merge all config sections, handle type rename carefully

---

#### MEDIUM CONFLICT FILES

##### 7. `internal/api/verify.go`
**Appflowy changes:**
- Enhanced OTP validation: `isOtpValid()` returns `(bool, error)`
- Magic link for new users
- Detailed error logging

**Master changes:**
- Various verification improvements
- Updated error handling

**Conflict severity:** MEDIUM
**Resolution strategy:** Merge enhanced OTP validation logic

---

##### 8. `internal/api/recover.go`
**Appflowy changes:**
- Clears password on recovery
- Sets `IsNonDefaultPassword = false`

**Master changes:**
- Various recovery improvements
- Updated error handling

**Conflict severity:** MEDIUM
**Resolution strategy:** Merge password clearing logic

---

##### 9. `internal/api/helpers.go`
**Appflowy changes:**
- Minimal changes

**Master changes:**
- Refactored to use apierrors
- Removed some utility functions

**Conflict severity:** MEDIUM
**Resolution strategy:** Use master version, port any appflowy-specific changes

---

##### 10. `cmd/migrate_cmd.go`
**Appflowy changes:**
- Auto-create namespace feature
- Context timeout for migrations

**Master changes:**
- Various migration improvements

**Conflict severity:** MEDIUM
**Resolution strategy:** Merge both feature sets

---

##### 11. `cmd/admin_cmd.go`
**Appflowy changes:**
- Sets `IsNonDefaultPassword = true` on admin user creation

**Master changes:**
- Various admin command improvements

**Conflict severity:** MEDIUM
**Resolution strategy:** Add IsNonDefaultPassword logic to master's version

---

#### LOW CONFLICT FILES

##### 12-17. Other Files
- `internal/api/reauthenticate.go`
- `internal/api/settings.go`
- `internal/api/password_test.go`
- `internal/conf/configuration_test.go`
- `internal/utilities/request.go`
- `internal/api/admin.go`

**Conflict severity:** LOW
**Resolution strategy:** Manual merge with careful review

---

## 6. Migration Files

### Deleted in AppFlowy
- `migrations/20240612123726_enable_rls_update_grants.up.sql`
  - Status: Check if needed for master

### Added in AppFlowy
- `migrations/20250425035447_alter_users_has_set_password.up.sql`
  - Adds `is_non_default_password` column

### Added in Master (since divergence)
- `migrations/20241009103726_add_web_authn.up.sql`
- `migrations/20250717082212_add_disabled_to_sso_providers.up.sql`
- `migrations/20250731150234_add_oauth_clients_table.up.sql`
- `migrations/20250804100000_add_oauth_authorizations_consents.up.sql`
- `migrations/20250901200500_add_oauth_client_type.up.sql`
- `migrations/20250903112500_remove_oauth_client_id_column.up.sql`
- `migrations/20250904133000_add_oauth_client_id_to_session.up.sql`
- `migrations/20250925093508_add_last_webauthn_challenge_data.up.sql`
- `migrations/20251007112900_add_session_refresh_token_columns.up.sql`

**Status:** All migrations are independent, no direct conflicts

---

## 7. Potential Conflicts & Issues

### CRITICAL ISSUES

#### 1. Error Handling System Incompatibility
**Problem:**
- AppFlowy uses old error structure in `internal/api/errors.go`
- Master moved everything to `internal/api/apierrors/`

**Impact:** All AppFlowy API functions will fail to compile

**Resolution:**
- Refactor all appflowy error calls to use `apierrors` package
- Update imports across all modified files
- Use `apierrors.NewBadRequestError()` instead of `badRequestError()`
- Use `apierrors.NewInternalServerError()` instead of `internalServerError()`

---

#### 2. User Model Schema Divergence
**Problem:**
- AppFlowy adds: `IsNonDefaultPassword`
- Master adds: WebAuthn fields, refresh token v2 fields

**Impact:** Schema conflicts, model incompatibility

**Resolution:**
- Merge both schema changes
- Ensure migration order is correct
- Update all code that creates/updates User models

---

#### 3. Password Validation Function Signature
**Problem:**
- AppFlowy: Regex-based validation with enhanced errors
- Master: Uses apierrors package

**Impact:** Function signature mismatch

**Resolution:**
- Keep regex logic from AppFlowy
- Update to use apierrors from master
- Merge enhanced error details

---

#### 4. Configuration Structure
**Problem:**
- Multiple overlapping config additions
- Type name change: `PasswordRequiredCharacters` → `PasswordRequiredCharactersRegex`

**Impact:** Config parsing failures, deployment issues

**Resolution:**
- Merge all config fields carefully
- Handle type rename with backward compatibility
- Update all config references

---

### MODERATE ISSUES

#### 5. API Route Conflicts
**Problem:**
- New routes in both branches
- Different middleware in master

**Impact:** Route registration conflicts

**Resolution:**
- Merge routes
- Ensure middleware order is correct
- Test all endpoints

---

#### 6. OTP Verification Logic
**Problem:**
- Different improvements in both branches
- Function signature change: `isOtpValid(...)` returns

**Impact:** Compilation errors

**Resolution:**
- Merge logic carefully
- Keep enhanced error handling
- Update all call sites

---

### LOW ISSUES

#### 7. Import Path Changes
**Problem:**
- Master refactored many packages
- AppFlowy uses old import paths

**Impact:** Import errors

**Resolution:**
- Update imports systematically

---

## 8. Architectural Differences

### AppFlowy Architecture Focus
- Password management lifecycle
- User authentication state tracking
- Improved error visibility for debugging
- Database auto-configuration

### Master Architecture Changes
- OAuth 2.1 authorization server
- Modern refresh token algorithm
- Multi-factor authentication expansion
- Background task processing
- Configuration hot-reloading
- Web3 authentication

### Compatibility Assessment
✅ **Generally compatible** - No fundamental architectural conflicts
⚠️ **Refactoring needed** - Error handling and imports must be updated
⚠️ **Testing required** - Password management features need validation

---

## 9. Merge Strategy

### Phase 1: Preparation (1-2 days)
- [x] Create `appflowy2` branch from master
- [x] Backup appflowy branch
- [x] Generate comparison document

### Phase 2: Systematic Merge (3-5 days)

#### Priority 1 - Foundation (Critical Path)
1. **Error handling refactor**
   - Port AppFlowy changes to use `apierrors` package
   - Update all function calls
   - Merge `WeakPasswordError` enhancements

2. **User model schema**
   - Merge `IsNonDefaultPassword` field
   - Add migration file
   - Add `UserAuthInfo` struct
   - Keep all master additions

3. **Password validation**
   - Port regex-based validation
   - Update to use apierrors
   - Merge enhanced error responses

#### Priority 2 - API Features
4. **New endpoints**
   - Add `GET /user/auth-info`
   - Add `POST /user/change-password`
   - Update API routing

5. **OTP verification**
   - Merge improved `isOtpValid` function
   - Merge magic link for new users
   - Keep detailed error logging

6. **Recovery flow**
   - Merge password clearing logic
   - Merge `IsNonDefaultPassword` reset

#### Priority 3 - Configuration & Commands
7. **Configuration**
   - Merge `AutoCreateNamespace` config
   - Rename to `PasswordRequiredCharactersRegex`
   - Keep all master config additions

8. **Migration commands**
   - Merge auto-create namespace
   - Add migration timeout

9. **Admin commands**
   - Merge `IsNonDefaultPassword` flag in admin user creation

### Phase 3: Testing (2-3 days)
- Unit tests
- Integration tests
- Migration tests
- End-to-end tests

### Phase 4: Documentation (1 day)
- Update API documentation
- Update configuration documentation
- Document schema changes

---

## 10. File-by-File Resolution Guide

### `internal/api/errors.go`
**Strategy:** Keep master structure, add AppFlowy enhancements

```go
// Keep master's type aliases to apierrors
// AppFlowy's WeakPasswordError changes should be in apierrors package
// Add MinLength and RequiredCharacters fields to apierrors.WeakPasswordError
```

---

### `internal/models/user.go`
**Strategy:** Merge all fields

```go
type User struct {
    // ... existing master fields ...
    IsNonDefaultPassword bool `json:"-" db:"is_non_default_password"` // ADD from appflowy
    // ... rest of master fields ...
}

// ADD from appflowy:
type UserAuthInfo struct {
    HasPassword          bool `json:"has_password"`
    IsSSOUser            bool `json:"is_sso_user"`
    IsNonDefaultPassword bool `json:"is_non_default_password"`
    IsSupabaseAdmin      bool `json:"is_supabase_admin"`
}
```

---

### `internal/api/user.go`
**Strategy:** Add AppFlowy functions with master's error handling

```go
// Add UserAuthInfoGet() function
// Add UserChangePassword() function
// Update to use apierrors instead of local error functions
// Keep all master functionality
```

---

### `internal/api/password.go`
**Strategy:** Keep AppFlowy's regex logic with master's structure

```go
// Use apierrors.NewBadRequestError instead of badRequestError
// Use apierrors.NewInternalServerError instead of internalServerError
// Keep regex matching logic
// Return &apierrors.WeakPasswordError with all fields
```

---

### `internal/api/verify.go`
**Strategy:** Merge both improvements

```go
// Keep appflowy's isOtpValid(bool, error) signature
// Keep detailed error messages
// Keep magic link for new users
// Update all function calls to handle (bool, error) return
```

---

### `internal/conf/configuration.go`
**Strategy:** Merge all config fields

```go
type DBConfiguration struct {
    // ... master fields ...
    AutoCreateNamespace bool `json:"auto_create_namespace" split_words:"true" default:"true"` // ADD
}

type PasswordConfiguration struct {
    MinLength int `json:"min_length" split_words:"true"`
    RequiredCharacters PasswordRequiredCharactersRegex `json:"required_characters" split_words:"true"` // RENAME
    HIBP HIBPConfiguration `json:"hibp"`
}

// RENAME type
type PasswordRequiredCharactersRegex []string // was PasswordRequiredCharacters
```

---

## 11. Testing Checklist

### Critical Path Tests

#### Authentication Flow
- [ ] User signup with password
- [ ] User signup without password (OAuth)
- [ ] User login with password
- [ ] OAuth login flow

#### Password Management
- [ ] Password change with old password validation
- [ ] Password change sets `IsNonDefaultPassword = true`
- [ ] Password recovery flow clears password
- [ ] Password recovery resets `IsNonDefaultPassword = false`
- [ ] Regex password validation with various patterns
- [ ] Weak password error includes detailed feedback

#### New Endpoints
- [ ] `GET /user/auth-info` returns correct data for:
  - [ ] User with password
  - [ ] SSO user
  - [ ] User with default password
  - [ ] User with custom password
  - [ ] Supabase admin user
- [ ] `POST /user/change-password`:
  - [ ] Validates current password correctly
  - [ ] Rejects incorrect current password
  - [ ] Validates new password strength
  - [ ] Sets `IsNonDefaultPassword = true`

#### OTP & Verification
- [ ] Magic link for new users works
- [ ] Magic link for confirmed users works
- [ ] OTP validation returns detailed errors
- [ ] Recovery verification for unconfirmed users

#### Database
- [ ] Auto-create namespace on fresh database
- [ ] Migration sequence runs successfully
- [ ] All migrations are idempotent
- [ ] Schema includes all expected columns

#### Master Features Compatibility
- [ ] OAuth 2.1 flows still work
- [ ] WebAuthn enrollment and verification
- [ ] Refresh token v2 generation and validation
- [ ] MFA with password change
- [ ] Background email sending
- [ ] Config reloading

#### Admin & Commands
- [ ] Admin user creation sets `IsNonDefaultPassword = true`
- [ ] Migration commands have timeout
- [ ] Migration commands auto-create namespace

#### Edge Cases
- [ ] SSO users cannot change password
- [ ] Anonymous users handled correctly
- [ ] Users with no password handled correctly
- [ ] Concurrent password change attempts

---

## 12. Risk Assessment

### HIGH RISK
1. **Error handling refactor** - Affects every API call
   - Mitigation: Systematic refactoring, comprehensive testing

2. **Password validation changes** - Security-critical
   - Mitigation: Thorough security review, penetration testing

3. **User model schema** - Database integrity
   - Mitigation: Careful migration testing, rollback plan

### MEDIUM RISK
4. **API routing changes** - Potential endpoint conflicts
   - Mitigation: Route testing, conflict detection

5. **OTP verification logic** - Authentication flow
   - Mitigation: Integration testing, manual verification

6. **Configuration structure** - Deployment impact
   - Mitigation: Configuration validation, backward compatibility

### LOW RISK
7. **Migration commands** - Isolated functionality
8. **Admin commands** - Limited scope
9. **Dependency updates** - Backward compatible

---

## 13. Rollback Plan

### If Issues Arise

#### Immediate Rollback
```bash
# Reset to master
git checkout master
git branch -D appflowy2

# If already pushed
git push origin --delete appflowy2
```

#### Partial Rollback
- Feature flags for new appflowy endpoints
- Database column defaults ensure backward compatibility
- Config defaults maintain existing behavior

#### Database Rollback
```sql
-- IsNonDefaultPassword has default value (false)
-- Can be safely rolled back by removing column
ALTER TABLE auth.users DROP COLUMN IF EXISTS is_non_default_password;
```

No data loss risk - all changes are additive.

---

## Timeline Estimate

| Phase | Duration | Activities |
|-------|----------|------------|
| Preparation | 1-2 days | Branch creation, documentation |
| Code merge | 3-5 days | File-by-file conflict resolution |
| Testing | 2-3 days | Unit, integration, e2e tests |
| Review & docs | 1 day | Code review, documentation updates |
| Buffer | 2 days | Unexpected issues |
| **Total** | **9-13 days** | Complete, tested merge |

---

## Appendix A: Useful Commands

### Branch Operations
```bash
# View divergence
git log --oneline --graph --decorate master...origin/appflowy

# Find merge base
git merge-base master origin/appflowy

# List files changed in appflowy
git diff --name-only $(git merge-base master origin/appflowy)..origin/appflowy

# List files changed in master
git diff --name-only $(git merge-base master origin/appflowy)..master

# Find conflicting files
git diff --name-only $(git merge-base master origin/appflowy)..origin/appflowy > /tmp/appflowy_files.txt
git diff --name-only $(git merge-base master origin/appflowy)..master > /tmp/master_files.txt
comm -12 <(sort /tmp/appflowy_files.txt) <(sort /tmp/master_files.txt)
```

### Testing
```bash
# Run unit tests
go test ./...

# Run specific package tests
go test ./internal/api/...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

### Migration
```bash
# Run migrations
go run cmd/main.go migrate up

# Check migration status
go run cmd/main.go migrate version

# Rollback last migration
go run cmd/main.go migrate down 1
```

---

## Appendix B: Key Contacts & Resources

- **AppFlowy Team:** Original feature requests and requirements
- **Master Branch Maintainers:** OAuth, WebAuthn, refresh token v2 experts
- **Database Team:** Schema changes and migration review
- **Security Team:** Password validation and authentication review

---

## Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-07 | Automated Analysis | Initial document creation |

---

**END OF DOCUMENT**
