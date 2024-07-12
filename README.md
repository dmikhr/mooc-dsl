MOOC DSL project aims to simplify the process of creating and managing educational tests across various platforms and systems.

Tutors often face the challenge of adapting their test descriptions to different Learning Management Systems (LMS) and standalone software. Each system typically has its own syntax and format, requiring tutors to invest significant time in understanding and translating their test descriptions.

This project addresses these challenges by introducing a Domain-Specific Language with simple and intuitive syntax for describing tests (see ```assets/sample.txt``` for DSL specifications). 

Tutors describe tests using DSL syntax, which are then processed and transformed into JSON format that can be integrated with third-party educational software, including LMS and custom educational platforms.

Project provides two main features:
* Syntax checker that outputs JSON with errors descriptions
* Parser which transforms DSL into JSON

Check ```assets``` folder for samples of JSON output and sample tests.
