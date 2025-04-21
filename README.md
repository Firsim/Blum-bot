# Blum-bot

Telegram channel for discussions and your questions: https://t.me/projectby

## Download:

The latest compiled version of the application for Windows is available for download on the releases page: https://github.com/Firsim/Blum-bot/releases

P.S. Versions for other operating systems will not be provided.

---

## Main Features of the Application:

1. **Automatic Clicker for Playing Blum Bot in Telegram**:
   - Recognition of game elements on the screen based on pixel color.
   - Automatic initiation of the next game.
   - Setting the number of games to be played.

2. **Configuration via `config.json` File**:
   - Color ranges for recognition (up to 10 colors).
   - Delays between clicks on the "Play" button.
   - Relative position of the "Play" button as percentages.

3. **Execution Control**:
   - Pause/Resume by pressing the 'P' key.
   - Limiting the number of games via the `-play` parameter.
   - Protection against multiple unintended clicks.

4. **Security**:
   - **IMPORTANT!** Administrator privileges are required for proper interaction with Windows system APIs.

---

### Example of Running via Command Line (cmd):

1. Open the "Start" menu and search for "Command Prompt" or "cmd".

2. Right-click on "Command Prompt" and select "Run as administrator".

3. In the opened window, navigate to the folder containing the program, for example:
   ```
   cd C:\Users\Username\Documents\Blum-bot
   ```

4. Run the program with the specified number of games (e.g., 5) using the `-play` flag:
   ```
   blum_bot.exe -play 5
   ```

   Where:
   - `blum_bot.exe` is the name of the program's executable file.
   - `-play 5` is the flag indicating the number of games (in this case, 5).

5. Wait for the startup message and follow the instructions in the console.

### Important Notes:
- If you encounter an error stating that administrator privileges are required, close the program and repeat steps 1–4, ensuring that the command prompt (`cmd`) is running with administrator rights.
