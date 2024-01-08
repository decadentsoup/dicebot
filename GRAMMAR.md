# Formula Grammar

This grammar is implemented by hand in `internal/lexer` and `internal/parser`.

```ebnf
(* Whitespace is ignored. Commas are considered whitespace. *)

    formula = equation*;
   equation = [name, "="], term;
       name = word, {word};

(* "e/md/as" is referring to the acronym "PEMDAS" for order of operations *)
       term = e term;
     e term = md term, {"^", md term};
    md term = as term, {("*" | "/"), as term};
    as term = factor, {("+" | "-"), factor};
bottom term = dice factor | unary term | "(", term, ")";
  dice term = [int], d, int;
 unary term = ["+" | "-"], int;

          d = "D" | "d";

        int = digit, {digit};
      digit = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9";

(* Word cannot match a dice term to avoid ambiguity. *)
       word = ("_" | letter), ("_" | letter | number);
     letter = ?any unicode letter?;
     number = ?any unicode number?;
```
