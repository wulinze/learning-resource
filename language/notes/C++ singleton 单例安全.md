#### C++ singleton 单例安全

##### 饿汉式和懒汉式

* 饿汉没有线程安全问题

~~~c++
//singleton.h
class singleton{
private:
    singleton(){}
    static singleton* p;
public:
    Singleton(const Singleton&)=delete;
    Singleton& operator=(const Singleton&)=delete;
    static singleton* instance();
    int a;
};

//singleton.cpp
singleton* singleton::p = new singleton();
singleton* singleton::instance()
{
    return p;
}

//main.c
int main(void)
{
    singleton * s = singleton::instance();
    s->a = 10;

    return 0;
}
~~~

* 懒汉的线程安全问题

  因为是非class static变量如果没有需要初始化

  **非线程安全版**

  ~~~c++
  //singleton.h
  class singleton{
  private:
      singleton(){}
      singleton* p;
  public:
      static singleton* instance();
      int a;
  };
  
  //singleton.cpp
  singleton* singleton::p = new singleton();
  singleton* singleton::instance()
  {
      if(p == nullptr){
          p = new singleton();
      }
      return p;
  }
  ~~~

  **线程安全(加锁版)**

  ~~~c++
  #include <iostream>
  #include <memory> // shared_ptr
  #include <mutex>  // mutex
  
  // version 2:
  // with problems below fixed:
  // 1. thread is safe now
  // 2. memory doesn't leak
  
  class Singleton{
  public:
      typedef std::shared_ptr<Singleton> Ptr;
      ~Singleton(){
          std::cout<<"destructor called!"<<std::endl;
      }
      Singleton(Singleton&)=delete;
      Singleton& operator=(const Singleton&)=delete;
      static Ptr get_instance(){
  
          // "double checked lock"
          if(m_instance_ptr==nullptr){
              std::lock_guard<std::mutex> lk(m_mutex);
              if(m_instance_ptr == nullptr){
                m_instance_ptr = std::shared_ptr<Singleton>(new Singleton);
              }
          }
          return m_instance_ptr;
      }
  
  
  private:
      Singleton(){
          std::cout<<"constructor called!"<<std::endl;
      }
      static Ptr m_instance_ptr;
      static std::mutex m_mutex;
  };
  
  // initialization static variables out of class
  Singleton::Ptr Singleton::m_instance_ptr = nullptr;
  std::mutex Singleton::m_mutex;
  
  int main(){
      Singleton::Ptr instance = Singleton::get_instance();
      Singleton::Ptr instance2 = Singleton::get_instance();
      return 0;
  }
  ~~~

  **线程安全(局部静态变量版)**

  ~~~c++
  #include <iostream>
  
  class Singleton
  {
  public:
      ~Singleton(){
          std::cout<<"destructor called!"<<std::endl;
      }
      Singleton(const Singleton&)=delete;
      Singleton& operator=(const Singleton&)=delete;
      static Singleton& get_instance(){
          static Singleton instance;
          return instance;
  
      }
  private:
      Singleton(){
          std::cout<<"constructor called!"<<std::endl;
      }
  };
  
  int main(int argc, char *argv[])
  {
      Singleton& instance_1 = Singleton::get_instance();
      Singleton& instance_2 = Singleton::get_instance();
      return 0;
  }
  ~~~

  最好不要把引用换乘指针进行操作，因为如果其中一个delete会导致内存泄漏

**单例模板**

* **CRTP 奇异递归模板的模式实现**

  基类模板的实现要点是：

  1. 构造函数需要是 **protected**，这样子类才能继承；
  2. 使用了奇异递归模板模式CRTP(Curiously recurring template pattern)
  3. get instance 方法和 2.2.3 的static local方法一个原理。
  4. 在这里基类的析构函数可以不需要 virtual ，因为子类在应用中只会用 Derived 类型，保证了析构时和构造时的类型一致

  ~~~c++
  // brief: a singleton base class offering an easy way to create singleton
  #include <iostream>
  
  template<typename T>
  class Singleton{
  public:
      static T& get_instance(){
          static T instance;
          return instance;
      }
      virtual ~Singleton(){
          std::cout<<"destructor called!"<<std::endl;
      }
      Singleton(const Singleton&)=delete;
      Singleton& operator =(const Singleton&)=delete;
  protected:
      Singleton(){
          std::cout<<"constructor called!"<<std::endl;
      }
  
  };
  /********************************************/
  // Example:
  // 1.friend class declaration is requiered!
  // 2.constructor should be private
  
  
  class DerivedSingle:public Singleton<DerivedSingle>{
     // !!!! attention!!!
     // needs to be friend in order to
     // access the private constructor/destructor
     friend class Singleton<DerivedSingle>;
  public:
     DerivedSingle(const DerivedSingle&)=delete;
     DerivedSingle& operator =(const DerivedSingle&)= delete;
  private:
     DerivedSingle()=default;
  };
  
  int main(int argc, char* argv[]){
      DerivedSingle& instance1 = DerivedSingle::get_instance();
      DerivedSingle& instance2 = DerivedSingle::get_instance();
      return 0;
  }
  ~~~

* **非友元声明**

  ~~~c++
  // brief: a singleton base class offering an easy way to create singleton
  #include <iostream>
  
  template<typename T>
  class Singleton{
  public:
      static T& get_instance() noexcept(std::is_nothrow_constructible<T>::value){
          static T instance{token()};
          return instance;
      }
      virtual ~Singleton() =default;
      Singleton(const Singleton&)=delete;
      Singleton& operator =(const Singleton&)=delete;
  protected:
      struct token{}; // helper class
      Singleton() noexcept=default;
  };
  
  
  /********************************************/
  // Example:
  // constructor should be public because protected `token` control the access
  
  
  class DerivedSingle:public Singleton<DerivedSingle>{
  public:
     DerivedSingle(token){
         std::cout<<"constructor called!"<<std::endl;
     }
  
     ~DerivedSingle(){
         std::cout<<"destructor called!"<<std::endl;
     }
     DerivedSingle(const DerivedSingle&)=delete;
     DerivedSingle& operator =(const DerivedSingle&)= delete;
  };
  
  int main(int argc, char* argv[]){
      DerivedSingle& instance1 = DerivedSingle::get_instance();
      DerivedSingle& instance2 = DerivedSingle::get_instance();
      return 0;
  }
  ~~~

  