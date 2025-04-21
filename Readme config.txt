 Description of Each Variable in the Configuration:

1. `button_text`:
   - Value: `"Play"`
   - Description: This variable is a remnant from the original version of the program and is only displayed in the console window.

2. `telegram_window_name`:
   - Value: `"Mini App: Blum"`
   - Description: The name of the Telegram window where the Blum mini-app runs. The program uses this name to locate and activate the correct window.

3. `colors`:
   - Value:
     ```json
     [
         {
             "red_range": 60,
             "green_range": 220,
             "blue_range": 10
         }
     ]
     ```
   - Description: Color ranges that the program uses to identify objects on the screen. In this case, one range is defined:
     - `red_range`: Red channel value
     - `green_range`: Green channel value
     - `blue_range`: Blue channel value
   - The program checks pixels on the screen and identifies those that fall within the specified color ranges.
   - Up to 10 colors can be defined in JSON format.

4. `click_delay_min`:
   - Value: `5`
   - Description: The minimum delay (in seconds) between clicks. The program will wait at least 5 seconds before performing the next click.

5. `click_delay_max`:
   - Value: `10`
   - Description: The maximum delay (in seconds) between clicks. The program will wait for a random duration between `click_delay_min` and `click_delay_max` before performing the next click.

6. `button_relative_position`:
   - Value:
     ```json
     {
         "x": 50,
         "y": 87.77
     }
     ```
   - Description: Relative coordinates of the button inside the Telegram window. The program uses these values to calculate the global coordinates of the button:
     - `x`: Horizontal offset of the button from the left edge of the window as a percentage.
     - `y`: Vertical offset of the button from the top edge of the window as a percentage.

---

 How It Works Together:

1. The program searches for the Telegram window with the name `"Mini App: Blum"` (`telegram_window_name`).
2. Inside this window, the program locates the "Play" button and calculates its global coordinates using the relative coordinates (`button_relative_position`).
3. The program performs clicks with a delay randomly chosen between `click_delay_min` and `click_delay_max` on the "Play" button.

---

 How Clicks Work Inside the Game:

- If the program finds a pixel with the following color:
  - Red channel = 60,
  - Green channel = 220,
  - Blue channel = 10,
  then this pixel is considered a match because it falls within the specified ranges (`colors`).

- The program clicks at this pixel's position relative to the top-left corner of the Telegram window.

