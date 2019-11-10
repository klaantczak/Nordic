#include <stdio.h>
#include "nordic32lf.h"

int main() {
  GoString name = {"src/nordic32/nordic32.json", 26};

  struct Init_return probs = Init(name);
  printf("Elements propbability of failure:");
  for (int i = 0; i < probs.r0; i++) {
    printf(" %f", ((double*)probs.r1)[i]);
  }
  printf("\n");

  int states[probs.r0];
  for (int i = 0; i < probs.r0; i++) {
    states[i] = 0;
  }

  for (int i = 40; i < 50; i++) {
    states[i] = 1;
  }

  struct Run_return res = Run(states);
  printf("Element states:");
  for (int i = 0; i < res.r0; i++) {
    printf(" %d", ((int*)res.r1)[i]);
  }
  printf("\n");
  return 0;
}
