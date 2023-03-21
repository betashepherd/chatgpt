import curses 
import random 

# 定义了方向
DIRS = {
    curses.KEY_UP: (-1, 0),
    curses.KEY_DOWN: (1, 0),
    curses.KEY_LEFT: (0, -1),
    curses.KEY_RIGHT: (0, 1)
}

# 初始化游戏场景
def init_scene():
    scene = []
    for r in range(20):
        scene.append(['.'] * 40)
    return scene

# 打印游戏场景
def print_scene(stdscr, scene, snake, food):
    stdscr.erase()
    for r in range(len(scene)):
        for c in range(len(scene[r])):
            stdscr.addstr(r, c, scene[r][c])
    for pos in snake:
        stdscr.addstr(pos[0], pos[1], '@')
    stdscr.addstr(food[0], food[1], '*')
    stdscr.refresh()

# 随机生成食物位置
def generate_food(scene, snake):
    while True:
        r = random.randint(0, len(scene)-1)
        c = random.randint(0, len(scene[0])-1)
        if scene[r][c] != '@':
            return (r, c)

# 判断蛇是否撞墙或者撞到自己
def check_collision(scene, snake):
    head = snake[0]
    if head[0] < 0 or head[0] >= len(scene) \
            or head[1] < 0 or head[1] >= len(scene[0]):
        return True
    if head in snake[1:]:
        return True
    return False

# 游戏循环
def game_loop(stdscr):
    # 初始化场景和蛇
    scene = init_scene()
    snake = [(9, 10), (9, 11), (9, 12)]
    direction = DIRS[curses.KEY_RIGHT]
    food = generate_food(scene, snake)

    # 游戏主循环
    while True:
        # 控制蛇移动速度
        stdscr.timeout(100)
        key = stdscr.getch()
        if key in DIRS:
            direction = DIRS[key]

        # 计算蛇头的下一个位置
        head = snake[0]
        head = (head[0] + direction[0], head[1] + direction[1])

        # 判断蛇是否吃到了食物
        if head == food:
            food = generate_food(scene, snake)
        else:
            tail = snake.pop()

            # 标记蛇之前的位置为空白
            scene[tail[0]][tail[1]] = '.'

        # 判断蛇是否撞墙或者撞到自己
        if check_collision(scene, snake):
            break

        # 向蛇的头部添加新位置，并标记此处为蛇的身体
        snake.insert(0, head)
        scene[head[0]][head[1]] = '@'

        # 打印游戏场景
        print_scene(stdscr, scene, snake, food)

    # 打印游戏结束信息
    stdscr.addstr(10, 20, 'Game Over!')
    stdscr.refresh()
    stdscr.getch()

# 初始化游戏
if __name__ == '__main__':
    curses.wrapper(game_loop)
