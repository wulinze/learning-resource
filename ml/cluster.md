### Clustering criterion

* 肘部法则-Elbow Method

  k-means聚类，将每个簇的质点与簇内样本点的平方距离误差和称为畸变程度。对于一个簇，它的畸变程度越低，代表簇内成员越紧密，畸变程度越高，代表簇内结构越松散。

  对于有一定区分度的数据，在达到某个临界点时畸变程度会得到极大改善，之后缓慢下降，这个临界点就可以考虑为聚类性能较好的点。

* 轮廓系数-Silhouette Coefficient
  $$
  a_i表示簇内系数\\
  b_{ij}表示簇间系数\\
  b_i = min\{{b_{ij}}\}\\
  s= {(b_i-a_i)}/{max(a_i, b_i)}
  $$
  
  s取值在-1和1之间，越靠近-1代表s更应该分到其他簇，越接近1，代表样本所在簇合理。
  
* Dunn Index(DI)

  * d(i,j)：簇间样本间距

  * $d^,(k)$：簇内样本距离
    $$
    D = \frac{\min_{1\leq i<j\leq n}d(i,j)}{\max_{1\leq k\leq n}d^,(k)}
    $$
    DI越大越好

