## Assembly conventions

The code in this library adheres to some code conventions.
It is advised to stick to these when using the code in this library.
That is, if you want your code to behave as expected.

### Register states

* The A, B and C registers are used to push arguments into a function.
* The A register will hold the return value from a function if applicable.
* When calling a subroutine, the A, B and C registers are considered free
  for all. Their contents are not preserved across calls. So if you are
  using them for anything important, you are required to save/restore them
  manually.
* The registers X, Y, Z, I and J are guaranteed to be preserved across
  function calls. This means that any function which anticipates
  usage of these registers must first push their value to the stack and
  pop them back once done. As show below:
  
    :my_super_duper_function
      ; Our function uses X and Y, so save them first.
      set push, x
      set push, y
      
      ; Do stuff with X and Y here...
      
      ; Once done, restore registers and exit.
      set y, pop
      set x, pop
      set pc, pop

