Output explore.gif

Set FontSize 30
Set Width 1500
Set Height 700

# Run the command by pressing enter.
Type " ./guntar explore $TARBALL -o ${OUTDIR}_explore"
Enter 1
Sleep 1s

# Go to first file
Down 2
Sleep 1s

# Go into file
Right 1
Sleep 2s

# Quit file
Left 1
Sleep 2s

# Select file
Type a
Sleep 1s

# Go up on a folder
Up 1
Sleep 1s
Enter 1
Sleep 1s

# Select a file 
Down 2
Sleep 1s
Type a
Sleep 1s

# Go back
Left 1
Sleep 1s

# Go back
Down 2
Type e
Sleep 3s

Ctrl+c
Ctrl+l
Type "listing all selected files!"
Sleep 1s
Ctrl+c
Ctrl+l

Type "ls -alR ${OUTDIR}_explore"
Sleep 1s
Enter  1
Sleep 5s

