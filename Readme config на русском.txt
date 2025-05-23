Описание каждой переменной в конфигурации:

1. `button_text`:
   - Значение: `"Играть"`
   - Описание: Эта переменная осталась из изначальной версии программы и отображается только в консоли окна.

2. `telegram_window_name`:
   - Значение: `"Mini App: Blum"`
   - Описание: Название окна Telegram, в котором работает мини-приложение Blum. Программа использует это имя для поиска и активации нужного окна.

3. `colors`:
   - Значение:
     ```json
     [
         {
             "red_range": 60,
             "green_range": 220,
             "blue_range": 10
         }
     ]
     ```
   - Описание: Диапазоны цветов, которые программа использует для поиска объектов на экране. В данном случае задан один диапазон:
     - `red_range`: Красный канал
     - `green_range`: Зелёный канал
     - `blue_range`: Синий канал
   - Программа проверяет пиксели на экране и ищет те, которые попадают в указанные диапазоны цветов.
   - Можно в формате json задать до 10 цветов

4. `click_delay_min`:
   - Значение: `5`
   - Описание: Минимальная задержка (в секундах) между кликами. Программа будет ждать не менее 5 секунд перед выполнением следующего клика.

5. `click_delay_max`:
   - Значение: `10`
   - Описание: Максимальная задержка (в секундах) между кликами. Программа будет ждать случайное время в интервале от `click_delay_min` до `click_delay_max` перед выполнением следующего клика.

6. `button_relative_position`:
   - Значение:
     ```json
     {
         "x": 50,
         "y": 87.77 
     }
     ```
   - Описание: Относительные координаты кнопки внутри окна Telegram. Программа использует эти значения для вычисления глобальных координат кнопки:
     - `x`: Горизонтальное смещение кнопки относительно левого края окна в процентах
     - `y`: Вертикальное смещение кнопки относительно верхнего края окна в процентах

---

 Как это работает вместе:

1. Программа ищет окно Telegram с именем `"Mini App: Blum"` (`telegram_window_name`).
2. Внутри этого окна программа ищет кнопку "Играть" и по координатам вычисляет её глобальные координаты, используя относительные координаты (`button_relative_position`.
3. Программа выполняет клики с задержкой, случайно выбранной в интервале от `click_delay_min` до `click_delay_max` по кнопке "Играть".

---

 Как происходят клики внутри игры:

- Если программа находит пиксель с цветом, где:
  - Красный канал = 60,
  - Зелёный канал = 220,
  - Синий канал = 10,
  то этот пиксель считается подходящим, так как он попадает в указанные диапазоны (`colors`).

- Программа кликает этому пикселю относительно левого верхнего угла окна Telegram.

