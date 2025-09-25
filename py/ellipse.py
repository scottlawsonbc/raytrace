import numpy as np
import matplotlib.pyplot as plt

dpi = 72
fig, ax1 = plt.subplots(1, 1, figsize=(600/dpi, 300/dpi), dpi=dpi)
t = np.linspace(0, 2*np.pi, 50000)
a = 1.0
b = 1.5
p = 6
c = np.cos(t)
s = np.sin(t)
x = np.abs(c)**(2/p) * np.sign(c) * a
y = np.abs(s)**(2/p) * np.sign(s) * b


import perlin
idealx = x.copy()
idealy = y.copy()

R = 512
shape = (R, R)
res = (4, 4)
noise = perlin.generate_perlin_noise_2d(shape, res)

# Always assume b is greater than or equal to a.
assert b >= a, "b must be greater than a"
# cx = (R-1)/2
# cy = (R-1)/2
cx = 0
cy = 0
X = (x  * (R-1) / b).astype(np.int)
Y = (y  * (R-1)/ b).astype(np.int)
d = np.sqrt((X-cx)**2 + (Y-cy)**2)
direction = np.array([X-cx,Y-cy]) / d
noise_component = direction * noise[X, Y] * 0.05
x += noise_component[0]
y += noise_component[1]
# plt.imshow(noise)
plt.plot(x, y, lw=2)
plt.plot(idealx, idealy, 'r', lw=2)
plt.gca().set_aspect('equal', adjustable='box')
plt.show()





ax1.plot(x, y)
ax1.axis('equal')
ax1.axis('off')
plt.show()
