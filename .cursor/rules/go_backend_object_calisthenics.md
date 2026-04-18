# Object Calisthenics - Style Guidelines

> **Priority: Simplicity > Rules**
> 
> If a rule makes code more complex, skip it. The goal is readable, maintainable code—not rule compliance.

## The 9 Rules

### 1. One level of indentation per method
- Max 3 levels. If more needed, extract to helper function.

### 2. Don't use `else`
- Use guard clauses. Happy path always to the left.

### 3. Wrap primitives
- Wrap primitives in dedicated types when they have domain meaning.
- Example: `type UserID string` — only if used in multiple places.

### 4. First class collections
- Collections with behavior deserve dedicated structs.

### 5. One dot per line (`obj.field.method()`)
- Multiple dots = coupling. Consider extracting to a variable or domain method.

### 6. Don't abbreviate
- Descriptive names: `userRepository`, not `ur`. Exception: `ctx`, `req`, `err` (Go convention).

### 7. Keep classes small
- Max 200-300 lines. If growing, split it.

### 8. Max 2 instance variables per struct
- More than 2 indicates multiple responsibilities. Consider splitting.

### 9. No getters/setters
- Group data with behavior. Avoid anemic domain models.

---

## When NOT to Force

| Scenario | Recommendation |
|----------|--------------|
| Struct with ~50 lines and 3+ fields | Keep, don't split just for rule 8 |
| Simple `GetX()` helper without logic | Keep getter if clear need |
| Simple if/else (2 branches) | Don't rewrite just for "no else" |
| Standard abbreviations (`ctx`, `err`, `req`) | Use freely |
| DTOs/input structs | May have getters for mapping |
| Code more readable without a rule | Skip the rule |

## Pragmatism Example

```go
// ✅ OK - don't split just because it has 3 fields
type User struct {
    ID    string
    Name  string
    Email string
}

// ❌ FORCED - creating type just for "rule 3" when string suffices
type UserID string
```

---

## Apply in Refactoring

Only apply when:
- Code is complex to read/maintain
- You or the team is having trouble understanding
- Adding new behavior

For new/simple code: write naturally and apply rules where it makes sense.