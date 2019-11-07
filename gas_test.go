package main

import "testing"

func TestAddGas(t *testing.T) {
	mixture := GasMixture{}
	mixture.Volume = 1000

	mixture.AddGas(GAS_O2, 50)

	//addition of new gas to mixture
	if amount := mixture.GetMolarValue(GAS_O2); amount != 50 {
		t.Errorf("Bad gas add. Added 50, got %f", amount)
	}
	if total := mixture.totalMolar; total != 50 {
		t.Errorf("Bad gas add. Added 50 to total, got %f", total)
	}

	//addition of gas to mixture that already contains some of the gas
	mixture.AddGas(GAS_O2, 100)
	if amount := mixture.GetMolarValue(GAS_O2); amount != 150 {
		t.Errorf("Bad gas add. Added 100 to 50, got %f", amount)
	}
	if total := mixture.totalMolar; total != 150 {
		t.Errorf("Bad gas add. Added 100 to total, got %f (expected 150)", total)
	}

	//addition of second gas to mixture, ensure it doesn't mess up the first
	mixture.AddGas(GAS_N2, 200)
	if amount := mixture.GetMolarValue(GAS_N2); amount != 200 {
		t.Errorf("Bad gas add. Added 200 N2, got N2 = %f", amount)
	}
	if amount := mixture.GetMolarValue(GAS_O2); amount != 150 {
		t.Errorf("Bad gas add. Added 200 N2, got O2 = %f", amount)
	}
	if total := mixture.totalMolar; total != 350 {
		t.Errorf("Bad gas add. Added 200 to total, got %f (expected 350)", total)
	}
}

func TestRemoveGas(t *testing.T) {
	mixture := GasMixture{}
	mixture.Volume = 1000

	mixture.AddGas(GAS_O2, 200)

	//remove some gas
	mixture.RemoveGas(GAS_O2, 100)
	if amount := mixture.GetMolarValue(GAS_O2); amount != 100 {
		t.Errorf("Bad gas remove. Removed 100 from 200, got %f", amount)
	}
	if total := mixture.totalMolar; total != 100 {
		t.Errorf("Bad gas remove. Removed 100 from total, got %f (expected 100)", total)
	}

	//attempt to remove gas that isn't there
	mixture.RemoveGas(GAS_N2, 100)
	if amount := mixture.GetMolarValue(GAS_O2); amount != 100 {
		t.Errorf("Bad gas remove. Removed N2 from O2 mixture, got O2 = %f", amount)
	}
	if total := mixture.totalMolar; total != 100 {
		t.Errorf("Bad gas remove. Removed none from total, got %f (expected 100)", total)
	}

	//remove exact amount of gas remaining
	mixture.RemoveGas(GAS_O2, 100)
	if amount := mixture.GetMolarValue(GAS_O2); amount != 0 {
		t.Errorf("Bad gas remove. Removed all O2, got O2 = %f", amount)
	}
	if total := mixture.totalMolar; total != 0 {
		t.Errorf("Bad gas remove. Removed all from total, got %f (expected 0)", total)
	}

	mixture.AddGas(GAS_O2, 200)

	//remove too much gas
	mixture.RemoveGas(GAS_O2, 300)
	if amount := mixture.GetMolarValue(GAS_O2); amount != 0 {
		t.Errorf("Bad gas remove. Removed 300 O2, got O2 = %f (expected 0)", amount)
	}
	if total := mixture.totalMolar; total != 0 {
		t.Errorf("Bad gas remove. Removed too much from total, got %f (expected 0)", total)
	}
}
