### 优化算法总结

#### 优化算法的整体公式

1. **当前梯度：$g_t=\nabla f(\omega_t)$** 
2. **根据梯度计算一阶动量$m_t$和二阶动量$V_t$**
3. **计算损失优化值：$\eta_t = \alpha \dot{} \frac{m_t}{\sqrt{V_t}}$   ($\alpha$代表学习率)**
4. **更新公式：$\omega_{t+1}=\omega_{t}-\eta_t$**

* **基本的GD：**

  $m_t = g_t$

  $V_t=I$

  $\eta_t = \alpha$ $\dot{}g_t$

  * BGD
  * SGD

* **GD with Momentum**

  $m_t = \beta_1 \dot{} m_{t-1}+(1-\beta_1)\dot{}g_t$

  使用累积梯度，不仅由当前梯度决定而且考虑了之前的梯度影响。

* **GD with Nesterov Acceleration**

  $g_t=\nabla f(\omega_t-\alpha\dot{}\frac{m_{t-1}}{\sqrt{V_t}})$ 

  当前梯度影响需要考虑到下一步的抉择。

* **AdaGrad**

  引入二阶动量

  $V_t=\sum_{\gamma=1}^tg_{\gamma}^2$

  随着训练样本增多学习率逐渐减小，减少但个样本的噪声，但是会阻止新的数据的学习

* **AdaDelta**

  通过加权进行求解，减少之前的影响增大当前样本的影响

  $V_t=\beta_2*V_{t-1}+(1-\beta_2)g_t^2$

* **Adam**

  $m_t = \beta_1 \dot{} m_{t-1}+(1-\beta_1)\dot{}g_t$

  $V_t=\beta_2*V_{t-1}+(1-\beta_2)g_t^2$

* **Nadam**

  Nesterov Acceleration+Adam
  
* **误差修正**：
  $$
  m_t = \frac{m_t}{1-\beta_1^t}\\
  V_t = \frac{V_t}{1-\beta_2^t}
  $$
  

