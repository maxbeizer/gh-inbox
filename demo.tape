# gh-inbox demo recording
# Run: vhs demo.tape

Output demo.gif

Set Shell zsh
Set FontFamily "JetBrainsMono Nerd Font"
Set FontSize 14
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set Padding 10
Set TypingSpeed 80ms

# Start gh-inbox in demo mode
Type "gh-inbox --demo"
Enter
Sleep 2s

# Navigate down through notifications
Sleep 500ms
Type "j"
Sleep 400ms
Type "j"
Sleep 400ms
Type "j"
Sleep 400ms

# Open preview panel
Type "p"
Sleep 2s

# Navigate with preview open
Type "j"
Sleep 1s
Type "j"
Sleep 1s

# Close preview
Type "p"
Sleep 1s

# Mark one as done
Type "e"
Sleep 1s

# Mark one as read
Type "I"
Sleep 1s

# Cycle filter to All
Type "f"
Sleep 1.5s

# Cycle back to Unread
Type "f"
Sleep 1s
Type "f"
Sleep 1s

# Search
Type "/"
Sleep 500ms
Type "api"
Enter
Sleep 1.5s

# Clear search
Escape
Sleep 1s

# Open help
Type "?"
Sleep 3s

# Close help
Type "?"
Sleep 1s

# Quit
Type "q"
Sleep 500ms
