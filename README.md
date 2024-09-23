https://threedots.tech/post/making-games-in-go/
https://www.youtube.com/live/UG0lghWfD-4?si=u483-otjTsKbR30y


Art style inspirations:
  - pixal art
  - caves of qud 
  - blasphemous
    
![caves of qud](https://www.cavesofqud.com/img/barathrums-study.png)
![blasphemous](https://static0.gamerantimages.com/wordpress/wp-content/uploads/2021/12/Blasphemous-CROPPED.jpg)


game idea

- a roguelike tag game
- timer till enemies get hostile
- obstacles to jump over
- obstacles to hide under
- can see only in light
- map is dark // fog of war // (memory mechanics from project zomboid where object remain on the place where you last saw them but dissapear slowly) (https://www.youtube.com/watch?v=nmE1OFPFLXc)
- when hostile should surviveand find the exit to go back one level
- learn powerup by defeating bosses
- bosses are hidden and not mandatory but enemies keeps getting harder, 2-3 levels after boss level enemies will have a chance to have the boss's abilities and the only way to get that abality is to go back a level and search for the boss


- Boss abalities will include stuff like teliportation, reducing players eyesight, AOE attacks, time manipulation, creating walls etc.

  
restructure ECS:

GEOMETRY
  position : X, Y
  sprite : image, hight, width
  boundingbox : polygon

PLAYER
  player : position, sprite

ENEMY
  enemy : position, sprite

LEVEL
  objects : position, sprite
  environment : []objects
  level: Environment, []enemy, player


game: level


current progress : 23/9/2024<br>
![23/9/2024](https://i.imgur.com/dU9rA6c.png)