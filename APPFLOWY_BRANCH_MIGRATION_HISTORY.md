# AppFlowy Branch Migration History

**Version:** 1.0
**Date:** 2025-11-08
**Status:** ✅ Migration Complete

---

## Quick Summary

The `appflowy` branch (28 commits with password management features) diverged from `master` 251 commits ago. We created `appflowy2` to systematically merge these features into current master.

**Result:** `appflowy2` = master (latest) + appflowy (features)

---

## Branch Relationships

```mermaid
gitGraph
    commit id: "10fa347 - Divergence Point (March 2025)"
    branch appflowy
    checkout appflowy
    commit id: "appflowy: 28 commits"
    commit id: "Password Management Features"
    commit id: "4377090 - tag: 0.8.0"

    checkout main
    commit id: "master: 251 commits"
    commit id: "OAuth 2.1, WebAuthn, RefreshV2, etc."
    commit id: "4e8275f - Current Master"

    branch appflowy2
    commit id: "Start from master"
    commit id: "Merge appflowy features (13 commits)"
    commit id: "60a5f95 - Ready for PR"

    checkout appflowy
    branch appflowy-backup
    commit id: "Backup created" type: HIGHLIGHT
```

---

## Branch Timeline

```mermaid
timeline
    title AppFlowy Branch Migration Timeline

    March 2025 : Divergence Point (10fa347)
               : appflowy branches from master

    March-Oct 2025 : appflowy Development
                   : 28 commits - password features
                   : Released as v0.8.0

    March-Nov 2025 : master Development
                   : 251 commits ahead
                   : OAuth2.1, WebAuthn, RefreshV2

    Nov 8, 2025 : Migration Day
                : Created appflowy2 from master
                : Created appflowy-backup
                : Merged 13 commits
                : Build ✅ Tests ⏳
```

---

## Why appflowy2 Exists

### The Problem

```mermaid
graph LR
    A[appflowy] -->|251 commits behind| B{Direct Merge?}
    B -->|❌ Massive Conflicts| C[17 files conflict]
    B -->|❌ Lose master work| D[OAuth, WebAuthn, etc.]
    B -->|❌ Rewrite history| E[Risky rebase]

    style B fill:#f99
    style C fill:#f99
    style D fill:#f99
    style E fill:#f99
```

### The Solution

```mermaid
graph TD
    M[master<br/>Latest] -->|git checkout -b| A2[appflowy2<br/>New Branch]
    AF[appflowy<br/>Features] -->|Manual Merge| A2
    A2 -->|13 commits| A2F[appflowy2<br/>Complete]
    A2F -->|PR| M2[master<br/>Updated]

    AF -->|Backup| AFB[appflowy-backup]

    style A2F fill:#9f9
    style M2 fill:#9f9
    style AFB fill:#99f
```

---

## Branch Status

| Branch | Status | Purpose | Delete After Merge? |
|--------|--------|---------|-------------------|
| `master` | ✅ Active | Main development | ❌ Never |
| `appflowy` | 🟡 Historical | Original work (v0.8.0) | ⚠️ After PR merged |
| `appflowy-backup` | 🔵 Backup | Preserves original | ✅ Optional |
| `appflowy2` | ✅ Active | Merge branch | ✅ After PR merged |

---

## Migration Flow

```mermaid
flowchart TD
    Start([Start Migration]) --> Prep[Phase 1: Preparation]
    Prep --> P1[Create appflowy2 from master]
    Prep --> P2[Create appflowy-backup]
    Prep --> P3[Generate analysis doc]

    P3 --> Merge[Phase 2: Merge Features]
    Merge --> M1[Enhance error handling]
    M1 --> M2[Update user model]
    M2 --> M3[Add migrations]
    M3 --> M4[Regex password validation]
    M4 --> M5[New API endpoints]
    M5 --> M6[Configuration changes]
    M6 --> M7[OTP improvements]
    M7 --> M8[Recovery flow]
    M8 --> M9[Migration commands]
    M9 --> M10[Admin commands]

    M10 --> Test[Phase 3: Test & Fix]
    Test --> T1[Build project]
    T1 --> T2{Build OK?}
    T2 -->|❌| Fix[Fix compilation errors]
    Fix --> T1
    T2 -->|✅| Done([Migration Complete])

    style Done fill:#9f9
    style Prep fill:#99f
    style Merge fill:#f96
    style Test fill:#ff9
```

---

## What Was Merged

```mermaid
mindmap
  root((appflowy2<br/>Features))
    Password Management
      IsNonDefaultPassword field
      User Auth Info API
      Password Change API
      Recovery flow reset
    Validation
      Regex patterns
      Enhanced errors
      Better feedback
    OTP
      Detailed errors
      Better logging
      Debug info
    Database
      Auto-create schema
      Migration timeout
      New migration file
    Configuration
      AutoCreateNamespace
      Regex config
```

---

## Commit Mapping Summary

```mermaid
graph LR
    subgraph appflowy[appflowy Branch - 17 commits]
        A1[Error handling]
        A2[User model]
        A3[API endpoints]
        A4[Password regex]
        A5[OTP errors]
        A6[Recovery]
        A7[Config]
        A8[Migration cmd]
    end

    subgraph appflowy2[appflowy2 Branch - 11 feature commits]
        B1[WeakPasswordError]
        B2[User model]
        B3[Migration file]
        B4[Regex validation]
        B5[API endpoints]
        B6[Config]
        B7[OTP improvements]
        B8[Recovery]
        B9[Migration cmd]
        B10[Admin cmd]
        B11[Fixes]
    end

    A1 --> B1
    A2 --> B2
    A2 --> B3
    A3 --> B5
    A4 --> B4
    A5 --> B7
    A6 --> B8
    A7 --> B6
    A7 --> B9
    A8 --> B9

    style appflowy fill:#e1f5ff
    style appflowy2 fill:#e8f5e9
```

**Mapping Types:**
- ✅ 1:1 (7 commits)
- 🔀 1:Many (2 commits split)
- 🔗 Many:1 (9 commits combined)

📄 **Detailed mapping:** See `APPFLOWY_COMMIT_MAPPING.md`

---

## Files Changed

```mermaid
pie title Files Modified by Category
    "API" : 6
    "Models" : 1
    "Config" : 1
    "Commands" : 2
    "Migrations" : 2
    "Errors" : 2
    "Docs" : 3
```

**Total:** 17 files changed, ~800 lines added

---

## Next Steps

```mermaid
stateDiagram-v2
    [*] --> RunTests: Current Status
    RunTests --> FixIssues: If failures
    RunTests --> PushBranch: If pass
    FixIssues --> RunTests
    PushBranch --> CreatePR
    CreatePR --> CodeReview
    CodeReview --> RequestChanges: If issues
    CodeReview --> Approve: If good
    RequestChanges --> MakeChanges
    MakeChanges --> CodeReview
    Approve --> MergePR
    MergePR --> DeleteBranches
    DeleteBranches --> [*]

    note right of RunTests
        go test ./...
    end note

    note right of PushBranch
        git push origin appflowy2
    end note

    note right of CreatePR
        Base: master
        Head: appflowy2
    end note
```

### Commands to Execute

```bash
# 1. Run tests
go test ./...

# 2. Push branch
git push origin appflowy2

# 3. Create PR
gh pr create --base master --head appflowy2 \
  --title "Merge AppFlowy password management features" \
  --body "$(cat APPFLOWY_MERGE_ANALYSIS.md)"

# 4. After merge
git checkout master
git pull origin master
git branch -d appflowy2
```

---

## Quick Reference

### View Branch Differences

```bash
# Compare appflowy to master
git log master..origin/appflowy --oneline

# Compare appflowy2 to master
git log master..appflowy2 --oneline

# See what appflowy2 merged
git diff master..appflowy2 --stat
```

### View Specific Changes

```bash
# Password validation
git show appflowy2:internal/api/password.go

# User model
git show appflowy2:internal/models/user.go

# API routes
git show appflowy2:internal/api/api.go
```

---

## Documentation Index

```mermaid
graph TD
    Start([Start Here]) --> This[BRANCH_MIGRATION_HISTORY.md<br/>Branch overview & timeline]

    This --> Detail{Need Details?}
    Detail -->|Commit Mapping| Mapping[COMMIT_MAPPING.md<br/>Exact 1:1 commit tracking]
    Detail -->|Technical Analysis| Analysis[MERGE_ANALYSIS.md<br/>600+ lines technical doc]
    Detail -->|Git History| Git[git log appflowy2]

    style This fill:#9cf
    style Mapping fill:#f9c
    style Analysis fill:#cf9
```

**Read in this order:**
1. **This document** - Understand why and what
2. **APPFLOWY_COMMIT_MAPPING.md** - Understand how (commit relationships)
3. **APPFLOWY_MERGE_ANALYSIS.md** - Understand details (technical deep-dive)

---

## Summary

```mermaid
graph LR
    A[appflowy<br/>28 commits<br/>Password Features]
    M[master<br/>251 commits<br/>OAuth, WebAuthn]

    A -->|Manual Merge| A2[appflowy2<br/>13 commits]
    M -->|Base| A2

    A2 -->|PR| M2[master<br/>All Features ✅]

    style A fill:#e1f5ff
    style M fill:#ffe1e1
    style A2 fill:#e8f5e9
    style M2 fill:#9f9
```

**Result:** Successfully merged 17 appflowy commits → 11 appflowy2 commits

**Status:** ✅ Build passes | ⏳ Tests pending | 📝 Ready for PR

---

**Last Updated:** 2025-11-08
**Document:** APPFLOWY_BRANCH_MIGRATION_HISTORY.md
**Version:** 1.0
