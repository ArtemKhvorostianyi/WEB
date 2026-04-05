import 'package:flutter/material.dart';

void main() {
  runApp(const MazutApp());
}

class MazutApp extends StatelessWidget {
  const MazutApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'Калькулятор мазуту',
      theme: ThemeData(
        useMaterial3: true,
        colorSchemeSeed: Colors.deepOrange,
      ),
      home: const MazutScreen(),
    );
  }
}

class MazutScreen extends StatefulWidget {
  const MazutScreen({super.key});

  @override
  State<MazutScreen> createState() => _MazutScreenState();
}

class _MazutScreenState extends State<MazutScreen> {

  final cg = TextEditingController(text: "92.40");
  final hg = TextEditingController(text: "1.83");
  final og = TextEditingController(text: "2.49");
  final sg = TextEditingController(text: "3.27");
  final qdaf = TextEditingController(text: "33.18");
  final wr = TextEditingController(text: "7.0");
  final ad = TextEditingController(text: "16.7");
  final vg = TextEditingController(text: "300");

  String? resultText;
  String? error;

  void calculate() {

    double Cg = double.tryParse(cg.text) ?? 0;
    double Hg = double.tryParse(hg.text) ?? 0;
    double Og = double.tryParse(og.text) ?? 0;
    double Sg = double.tryParse(sg.text) ?? 0;
    double Qdaf = double.tryParse(qdaf.text) ?? 0;
    double Wr = double.tryParse(wr.text) ?? 0;
    double Ad = double.tryParse(ad.text) ?? 0;
    double Vg = double.tryParse(vg.text) ?? 0;

    if (Wr + Ad >= 100) {
      setState(() {
        error = "Помилка: W + A ≥ 100%";
        resultText = null;
      });
      return;
    }

    if ((Cg + Hg + Og + Sg - 100).abs() > 0.5) {
      setState(() {
        error = "Сума Cг + Hг + Oг + Sг повинна ≈ 100%";
        resultText = null;
      });
      return;
    }

    double K = (100 - Wr - Ad) / 100;

    double Cr = Cg * K;
    double Hr = Hg * K;
    double Or = Og * K;
    double Sr = Sg * K;
    double Vr = Vg * K;

    double Qr = Qdaf * K - 0.025 * Wr;

    setState(() {
      error = null;

      resultText = """
Для складу горючої маси мазуту, що задано наступними параметрами:

Hг = ${Hg.toStringAsFixed(2)} %;
Cг = ${Cg.toStringAsFixed(2)} %;
Sг = ${Sg.toStringAsFixed(2)} %;
Oг = ${Og.toStringAsFixed(2)} %;
Vг = ${Vg.toStringAsFixed(1)} мг/кг;
Wг = ${Wr.toStringAsFixed(2)} %;
Aг = ${Ad.toStringAsFixed(2)} %;

та нижчою теплотою згоряння горючої маси мазуту
Qi daf = ${Qdaf.toStringAsFixed(2)} МДж/кг,

Склад робочої маси мазуту становитиме:

Hр = ${Hr.toStringAsFixed(2)} %;
Cр = ${Cr.toStringAsFixed(2)} %;
Sр = ${Sr.toStringAsFixed(2)} %;
Oр = ${Or.toStringAsFixed(2)} %;
Vр = ${Vr.toStringAsFixed(1)} мг/кг;
Aр = ${Ad.toStringAsFixed(2)} %;

Нижча теплота згоряння мазуту на робочу масу
за заданим складом компонентів палива становить:

Qr = ${Qr.toStringAsFixed(2)} МДж/кг.
""";
    });
  }

  Widget inputField(TextEditingController controller, String label) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: TextField(
        controller: controller,
        decoration: InputDecoration(
          labelText: label,
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
        keyboardType:
            const TextInputType.numberWithOptions(decimal: true),
      ),
    );
  }

  Widget resultCard() {
    if (resultText == null) return const SizedBox();

    return Card(
      elevation: 4,
      margin: const EdgeInsets.only(top: 20),
      shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Text(
          resultText!,
          style: const TextStyle(
            fontSize: 16,
            height: 1.5,
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {

    return Scaffold(
      appBar: AppBar(
        title: const Text("Калькулятор мазуту"),
        centerTitle: true,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [

            inputField(cg, "Вуглець Cг, %"),
            inputField(hg, "Водень Hг, %"),
            inputField(og, "Кисень Oг, %"),
            inputField(sg, "Сірка Sг, %"),
            inputField(qdaf, "Qi daf, МДж/кг"),
            inputField(wr, "Вологість Wг, %"),
            inputField(ad, "Зольність Aг, %"),
            inputField(vg, "Ванадій Vг, мг/кг"),

            const SizedBox(height: 20),

            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: calculate,
                style: ElevatedButton.styleFrom(
                  padding:
                      const EdgeInsets.symmetric(vertical: 15),
                ),
                child: const Text(
                  "Обчислити",
                  style: TextStyle(fontSize: 18),
                ),
              ),
            ),

            if (error != null)
              Padding(
                padding: const EdgeInsets.only(top: 15),
                child: Text(
                  error!,
                  style: const TextStyle(
                      color: Colors.red,
                      fontSize: 16),
                ),
              ),

            resultCard(),

          ],
        ),
      ),
    );
  }
}
